package test

import (
	"context"
	"os"
	"testing"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/web/mgmt/web"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tsalright/terratest/modules/azure"
)

const (
	// AzureSubscriptionID is an optional env variable supported by the `azurerm` Terraform provider to
	// designate a target Azure subscription ID
	AzureSubscriptionID = "ARM_SUBSCRIPTION_ID"

	// AzureResGroupName is an optional env variable custom to Terratest to designate a target Azure resource group
	AzureResGroupName = "tfmodulevalidation-test-group"
)

type planValidationArgs struct {
	resGroupName        string
	subID               string
	expectedPlanName    string
	expectedSkuSize     string
	expectedSkuTier     string
	expectedSkuCapacity int32
	expectedKind        string
	expectedReserved    bool
}

func TestTerraformBasicExample(t *testing.T) {
	// Arrange
	terraformOptions := &terraform.Options{
		TerraformDir: "../example/basic/.",
	}
	// Leaving subID empty to force it to use the environment variable that Terraform requires to run in CI
	planArgs := planValidationArgs{
		resGroupName:        "tfmodulevalidation-test-group",
		subID:               "",
		expectedPlanName:    "basicPlanSample-test-sharedplan-0-westus2",
		expectedSkuSize:     "S1",
		expectedSkuTier:     "Standard",
		expectedSkuCapacity: 1,
		expectedKind:        "Windows",
		expectedReserved:    false,
	}

	defer terraform.Destroy(t, terraformOptions)

	// Act
	terraform.InitAndApply(t, terraformOptions)

	// Assert
	validatePlanContent(t, terraformOptions, &planArgs)
}

func TestTerraformConsumptionExample(t *testing.T) {
	// Arrange
	terraformOptions := &terraform.Options{
		TerraformDir: "../example/consumption/.",
	}
	// Leaving subID empty to force it to use the environment variable that Terraform requires to run in CI
	planArgs := planValidationArgs{
		resGroupName:        AzureResGroupName,
		subID:               "",
		expectedPlanName:    "consumptionPlanSample-test-sharedplan-0-westus2",
		expectedSkuSize:     "Y1",
		expectedSkuTier:     "Dynamic",
		expectedSkuCapacity: 0,
		expectedKind:        "functionapp",
		expectedReserved:    false,
	}

	defer terraform.Destroy(t, terraformOptions)

	// Act
	terraform.InitAndApply(t, terraformOptions)

	validatePlanContent(t, terraformOptions, &planArgs)
}

func TestTerraformLinuxExample(t *testing.T) {
	// Arrange
	terraformOptions := &terraform.Options{
		TerraformDir: "../example/linux/.",
	}
	// Leaving subID empty to force it to use the environment variable that Terraform requires to run in CI
	planArgs := planValidationArgs{
		resGroupName:        "tfmodulevalidation-test-group",
		subID:               "",
		expectedPlanName:    "linuxPlanSample-test-sharedplan-0-westus2",
		expectedSkuSize:     "S1",
		expectedSkuTier:     "Standard",
		expectedSkuCapacity: 1,
		expectedKind:        "linux",
		expectedReserved:    true,
	}

	defer terraform.Destroy(t, terraformOptions)

	// Act
	terraform.InitAndApply(t, terraformOptions)

	// Assert
	validatePlanContent(t, terraformOptions, &planArgs)
}

func validatePlanContent(t *testing.T, terraformOptions *terraform.Options, args *planValidationArgs) {
	assert := assert.New(t)

	outputValue := terraform.Output(t, terraformOptions, "sharedplanids")
	assert.NotNil(outputValue)
	assert.Contains(outputValue, args.expectedPlanName)

	plan := GetAppServicePlan(t, args.expectedPlanName)
	if args.expectedSkuSize != "" {
		assert.Equal(args.expectedSkuSize, *plan.Sku.Size)
	}

	if args.expectedSkuTier != "" {
		assert.Equal(args.expectedSkuTier, *plan.Sku.Tier)
	}

	if args.expectedSkuCapacity != 0 {
		assert.Equal(args.expectedSkuCapacity, *plan.Sku.Capacity)
	}

	if args.expectedKind != "" {
		assert.Equal(args.expectedKind, *plan.Kind)
	}

	assert.Equal(args.expectedReserved, *plan.Reserved)
}

func GetAppServicePlan(t *testing.T, planName string) *web.AppServicePlan {
	plan, err := getAppServicePlanE(planName)
	require.NoError(t, err)

	return plan
}

func getAppServicePlanE(planName string) (*web.AppServicePlan, error) {
	client, err := getAppServicePlanClient()
	if err != nil {
		return nil, err
	}

	plan, err := client.Get(context.Background(), AzureResGroupName, planName)
	if err != nil {
		return nil, err
	}

	return &plan, nil
}

func getAppServicePlanClient() (*web.AppServicePlansClient, error) {
	subID := os.Getenv(AzureSubscriptionID)

	// Create an AppServicePlanClient
	planClient := web.NewAppServicePlansClient(subID)

	authorizer, err := azure.NewAuthorizer()
	if err != nil {
		return nil, err
	}

	planClient.Authorizer = *authorizer

	return &planClient, nil
}
