package test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/web/mgmt/web"
	"github.com/Azure/azure-sdk-for-go/services/preview/monitor/mgmt/2018-11-01-preview/insights"
	"github.com/gruntwork-io/terratest/modules/azure"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	// AzureSubscriptionID is an optional env variable supported by the `azurerm` Terraform provider to
	// designate a target Azure subscription ID
	AzureSubscriptionID = "ARM_SUBSCRIPTION_ID"

	// AzureResGroupName is an optional env variable custom to Terratest to designate a target Azure resource group
	AzureResGroupName = "tfmodulevalidation-test-group"
)

type planValidationArgs struct {
	expectedPlanName    string
	expectedSkuSize     string
	expectedSkuTier     string
	expectedSkuCapacity int32
	expectedKind        string
	expectedReserved    bool
}

type metricAlertsValidationArgs struct {
	expectedRuleName          string
	expectedMetricNameSpace   string
	expectedCPUMetricName     string
	expectedCPUOperator       string
	expectedCPUAggregation    string
	expectedCPUThreshold      float64
	expectedDiskMetricName    string
	expectedDiskOperator      string
	expectedDiskAggregation   string
	expectedDiskThreshold     float64
	expectedMemoryMetricName  string
	expectedMemoryOperator    string
	expectedMemoryAggregation string
	expectedMemoryThreshold   float64
	expectedHTTPMetricName    string
	expectedHTTPOperator      string
	expectedHTTPAggregation   string
	expectedHTTPThreshold     float64
}

type metricTriggerValidationArgs struct {
	metricName      string
	timeGrain       string
	statistic       string
	timeWindow      string
	timeAggregation string
	operator        string
	threshold       float64
}

type scaleActionValidationArgs struct {
	direction  string
	actionType string
	value      string
	cooldown   string
}

type ruleValidationArgs struct {
	metricTrigger *metricTriggerValidationArgs
	scaleAction   *scaleActionValidationArgs
}

type profileValidationArgs struct {
	name            string
	defaultCapacity string
	minimumCapacity string
	maximumCapacity string
	rules           *[]ruleValidationArgs
}

func TestTerraformBasicExample(t *testing.T) {
	// Arrange
	terraformOptions := &terraform.Options{
		TerraformDir: "../example/basic/.",
	}

	planArgs := planValidationArgs{
		expectedPlanName:    "basicPlanSample-test-sharedplan-0-westus2",
		expectedSkuSize:     "S1",
		expectedSkuTier:     "Standard",
		expectedSkuCapacity: 1,
		expectedKind:        "Windows",
		expectedReserved:    false,
	}

	metricAlertsArgs := metricAlertsValidationArgs{
		expectedCPUAggregation:    "Average",
		expectedCPUMetricName:     "CpuPercentage",
		expectedCPUOperator:       "GreaterThan",
		expectedCPUThreshold:      70,
		expectedDiskAggregation:   "Average",
		expectedDiskMetricName:    "DiskQueueLength",
		expectedDiskOperator:      "GreaterThan",
		expectedDiskThreshold:     100,
		expectedHTTPAggregation:   "Average",
		expectedHTTPMetricName:    "HttpQueueLength",
		expectedHTTPOperator:      "GreaterThan",
		expectedHTTPThreshold:     100,
		expectedMemoryAggregation: "Average",
		expectedMemoryMetricName:  "MemoryPercentage",
		expectedMemoryOperator:    "GreaterThan",
		expectedMemoryThreshold:   90,
		expectedMetricNameSpace:   "Microsoft.Web/serverfarms",
		expectedRuleName:          "basicPlanSample-test-sharedplan-alerts",
	}

	defer terraform.Destroy(t, terraformOptions)

	// Act
	terraform.InitAndApply(t, terraformOptions)

	// Assert
	validatePlanContent(t, terraformOptions, &planArgs)
	validateMetricAlertsContent(t, &metricAlertsArgs)
	validateAutoscaleSettingsContent(t, planArgs.expectedPlanName)
}

func TestTerraformConsumptionExample(t *testing.T) {
	// Arrange
	terraformOptions := &terraform.Options{
		TerraformDir: "../example/consumption/.",
	}

	planArgs := planValidationArgs{
		expectedPlanName:    "consumptionPlanSample-test-sharedplan-0-westus2",
		expectedSkuSize:     "Y1",
		expectedSkuTier:     "Dynamic",
		expectedSkuCapacity: 0,
		expectedKind:        "functionapp",
		expectedReserved:    false,
	}

	metricAlertsArgs := metricAlertsValidationArgs{
		expectedCPUAggregation:    "Average",
		expectedCPUMetricName:     "CpuPercentage",
		expectedCPUOperator:       "GreaterThan",
		expectedCPUThreshold:      70,
		expectedDiskAggregation:   "Average",
		expectedDiskMetricName:    "DiskQueueLength",
		expectedDiskOperator:      "GreaterThan",
		expectedDiskThreshold:     100,
		expectedHTTPAggregation:   "Average",
		expectedHTTPMetricName:    "HttpQueueLength",
		expectedHTTPOperator:      "GreaterThan",
		expectedHTTPThreshold:     100,
		expectedMemoryAggregation: "Average",
		expectedMemoryMetricName:  "MemoryPercentage",
		expectedMemoryOperator:    "GreaterThan",
		expectedMemoryThreshold:   90,
		expectedMetricNameSpace:   "Microsoft.Web/serverfarms",
		expectedRuleName:          "consumptionPlanSample-test-sharedplan-alerts",
	}

	defer terraform.Destroy(t, terraformOptions)

	// Act
	terraform.InitAndApply(t, terraformOptions)

	// Assert
	validatePlanContent(t, terraformOptions, &planArgs)
	validateMetricAlertsContent(t, &metricAlertsArgs)
}

func TestTerraformLinuxExample(t *testing.T) {
	// Arrange
	terraformOptions := &terraform.Options{
		TerraformDir: "../example/linux/.",
	}

	planArgs := planValidationArgs{
		expectedPlanName:    "linuxPlanSample-test-sharedplan-0-westus2",
		expectedSkuSize:     "S1",
		expectedSkuTier:     "Standard",
		expectedSkuCapacity: 1,
		expectedKind:        "linux",
		expectedReserved:    true,
	}

	metricAlertsArgs := metricAlertsValidationArgs{
		expectedCPUAggregation:    "Average",
		expectedCPUMetricName:     "CpuPercentage",
		expectedCPUOperator:       "GreaterThan",
		expectedCPUThreshold:      70,
		expectedDiskAggregation:   "Average",
		expectedDiskMetricName:    "DiskQueueLength",
		expectedDiskOperator:      "GreaterThan",
		expectedDiskThreshold:     100,
		expectedHTTPAggregation:   "Average",
		expectedHTTPMetricName:    "HttpQueueLength",
		expectedHTTPOperator:      "GreaterThan",
		expectedHTTPThreshold:     100,
		expectedMemoryAggregation: "Average",
		expectedMemoryMetricName:  "MemoryPercentage",
		expectedMemoryOperator:    "GreaterThan",
		expectedMemoryThreshold:   90,
		expectedMetricNameSpace:   "Microsoft.Web/serverfarms",
		expectedRuleName:          "linuxPlanSample-test-sharedplan-alerts",
	}

	defer terraform.Destroy(t, terraformOptions)

	// Act
	terraform.InitAndApply(t, terraformOptions)

	// Assert
	validatePlanContent(t, terraformOptions, &planArgs)
	validateMetricAlertsContent(t, &metricAlertsArgs)
	validateAutoscaleSettingsContent(t, planArgs.expectedPlanName)
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

func validateMetricAlertsContent(t *testing.T, args *metricAlertsValidationArgs) {
	assert := assert.New(t)

	metricAlerts := GetMetricAlertsResource(t, args.expectedRuleName)

	criteriaList, _ := metricAlerts.Criteria.AsMetricAlertSingleResourceMultipleMetricCriteria()

	for _, criteria := range *criteriaList.AllOf {
		fmt.Print(*criteria.MetricName)
		assert.Equal(args.expectedMetricNameSpace, *criteria.MetricNamespace)

		if args.expectedCPUMetricName == *criteria.MetricName {
			assert.Equal(args.expectedCPUOperator, string(criteria.Operator))
			assert.Equal(args.expectedCPUThreshold, *criteria.Threshold)
			assert.Equal(args.expectedCPUAggregation, criteria.TimeAggregation)
		}

		if args.expectedDiskMetricName == *criteria.MetricName {
			assert.Equal(args.expectedDiskOperator, string(criteria.Operator))
			assert.Equal(args.expectedDiskThreshold, *criteria.Threshold)
			assert.Equal(args.expectedDiskAggregation, criteria.TimeAggregation)
		}

		if args.expectedMemoryMetricName == *criteria.MetricName {
			assert.Equal(args.expectedMemoryOperator, string(criteria.Operator))
			assert.Equal(args.expectedMemoryThreshold, *criteria.Threshold)
			assert.Equal(args.expectedMemoryAggregation, criteria.TimeAggregation)
		}

		if args.expectedHTTPMetricName == *criteria.MetricName {
			assert.Equal(args.expectedHTTPOperator, string(criteria.Operator))
			assert.Equal(args.expectedHTTPThreshold, *criteria.Threshold)
			assert.Equal(args.expectedHTTPAggregation, criteria.TimeAggregation)
		}
	}
}

func validateAutoscaleSettingsContent(t *testing.T, resourceName string) {
	assert := assert.New(t)

	autoscaleExpected := buildAutoscaleValidationData()
	autoscaleSettings := GetAutoscaleSettingsResource(t, resourceName)
	for _, profile := range *autoscaleSettings.Profiles {
		assert.Equal(autoscaleExpected.name, *profile.Name)
		assert.Equal(autoscaleExpected.defaultCapacity, *profile.Capacity.Default)
		assert.Equal(autoscaleExpected.minimumCapacity, *profile.Capacity.Minimum)
		assert.Equal(autoscaleExpected.maximumCapacity, *profile.Capacity.Maximum)

		assert.Equal(len(*autoscaleExpected.rules), len(*profile.Rules))
		for _, rule := range *profile.Rules {
			assert.True(isRulePresent(&rule, autoscaleExpected.rules))
		}
	}
}

func isRulePresent(rule *insights.ScaleRule, expectedRules *[]ruleValidationArgs) bool {
	for _, expectedRule := range *expectedRules {
		if *rule.MetricTrigger.MetricName == expectedRule.metricTrigger.metricName &&
			*rule.MetricTrigger.TimeGrain == expectedRule.metricTrigger.timeGrain &&
			string(rule.MetricTrigger.Statistic) == expectedRule.metricTrigger.statistic &&
			*rule.MetricTrigger.TimeWindow == expectedRule.metricTrigger.timeWindow &&
			string(rule.MetricTrigger.TimeAggregation) == expectedRule.metricTrigger.timeAggregation &&
			string(rule.MetricTrigger.Operator) == expectedRule.metricTrigger.operator &&
			*rule.MetricTrigger.Threshold == expectedRule.metricTrigger.threshold &&
			string(rule.ScaleAction.Direction) == expectedRule.scaleAction.direction &&
			string(rule.ScaleAction.Type) == expectedRule.scaleAction.actionType &&
			*rule.ScaleAction.Value == expectedRule.scaleAction.value &&
			*rule.ScaleAction.Cooldown == expectedRule.scaleAction.cooldown {
			return true
		}
	}
	return false
}

func buildAutoscaleValidationData() *profileValidationArgs {
	profile := profileValidationArgs{
		name:            "defaultProfile",
		defaultCapacity: "1",
		maximumCapacity: "10",
		minimumCapacity: "1",
	}

	profile.rules = &[]ruleValidationArgs{
		{
			metricTrigger: &metricTriggerValidationArgs{
				metricName:      "CpuPercentage",
				timeGrain:       "PT1M",
				statistic:       "Average",
				timeWindow:      "PT5M",
				timeAggregation: "Average",
				operator:        "GreaterThan",
				threshold:       90,
			},
			scaleAction: &scaleActionValidationArgs{
				direction:  "Increase",
				actionType: "ChangeCount",
				value:      "2",
				cooldown:   "PT5M",
			},
		},
		{
			metricTrigger: &metricTriggerValidationArgs{
				metricName:      "CpuPercentage",
				timeGrain:       "PT1M",
				statistic:       "Average",
				timeWindow:      "PT5M",
				timeAggregation: "Average",
				operator:        "GreaterThan",
				threshold:       75,
			},
			scaleAction: &scaleActionValidationArgs{
				direction:  "Increase",
				actionType: "ChangeCount",
				value:      "1",
				cooldown:   "PT5M",
			},
		},
		{
			metricTrigger: &metricTriggerValidationArgs{
				metricName:      "CpuPercentage",
				timeGrain:       "PT1M",
				statistic:       "Average",
				timeWindow:      "PT15M",
				timeAggregation: "Average",
				operator:        "LessThanOrEqual",
				threshold:       50,
			},
			scaleAction: &scaleActionValidationArgs{
				direction:  "Decrease",
				actionType: "ChangeCount",
				value:      "1",
				cooldown:   "PT15M",
			},
		},
		{
			metricTrigger: &metricTriggerValidationArgs{
				metricName:      "HttpQueueLength",
				timeGrain:       "PT1M",
				statistic:       "Average",
				timeWindow:      "PT5M",
				timeAggregation: "Average",
				operator:        "GreaterThan",
				threshold:       100,
			},
			scaleAction: &scaleActionValidationArgs{
				direction:  "Increase",
				actionType: "ChangeCount",
				value:      "1",
				cooldown:   "PT5M",
			},
		},
		{
			metricTrigger: &metricTriggerValidationArgs{
				metricName:      "HttpQueueLength",
				timeGrain:       "PT1M",
				statistic:       "Average",
				timeWindow:      "PT15M",
				timeAggregation: "Average",
				operator:        "LessThanOrEqual",
				threshold:       50,
			},
			scaleAction: &scaleActionValidationArgs{
				direction:  "Decrease",
				actionType: "ChangeCount",
				value:      "1",
				cooldown:   "PT15M",
			},
		},
		{
			metricTrigger: &metricTriggerValidationArgs{
				metricName:      "HttpQueueLength",
				timeGrain:       "PT1M",
				statistic:       "Average",
				timeWindow:      "PT5M",
				timeAggregation: "Average",
				operator:        "GreaterThan",
				threshold:       200,
			},
			scaleAction: &scaleActionValidationArgs{
				direction:  "Increase",
				actionType: "ChangeCount",
				value:      "2",
				cooldown:   "PT5M",
			},
		},
		{
			metricTrigger: &metricTriggerValidationArgs{
				metricName:      "MemoryPercentage",
				timeGrain:       "PT1M",
				statistic:       "Average",
				timeWindow:      "PT5M",
				timeAggregation: "Average",
				operator:        "GreaterThan",
				threshold:       85,
			},
			scaleAction: &scaleActionValidationArgs{
				direction:  "Increase",
				actionType: "ChangeCount",
				value:      "1",
				cooldown:   "PT5M",
			},
		},
		{
			metricTrigger: &metricTriggerValidationArgs{
				metricName:      "MemoryPercentage",
				timeGrain:       "PT1M",
				statistic:       "Average",
				timeWindow:      "PT15M",
				timeAggregation: "Average",
				operator:        "LessThanOrEqual",
				threshold:       65,
			},
			scaleAction: &scaleActionValidationArgs{
				direction:  "Decrease",
				actionType: "ChangeCount",
				value:      "1",
				cooldown:   "PT15M",
			},
		},
		{
			metricTrigger: &metricTriggerValidationArgs{
				metricName:      "DiskQueueLength",
				timeGrain:       "PT1M",
				statistic:       "Average",
				timeWindow:      "PT5M",
				timeAggregation: "Average",
				operator:        "GreaterThan",
				threshold:       100,
			},
			scaleAction: &scaleActionValidationArgs{
				direction:  "Increase",
				actionType: "ChangeCount",
				value:      "1",
				cooldown:   "PT5M",
			},
		},
		{
			metricTrigger: &metricTriggerValidationArgs{
				metricName:      "DiskQueueLength",
				timeGrain:       "PT1M",
				statistic:       "Average",
				timeWindow:      "PT15M",
				timeAggregation: "Average",
				operator:        "LessThanOrEqual",
				threshold:       50,
			},
			scaleAction: &scaleActionValidationArgs{
				direction:  "Decrease",
				actionType: "ChangeCount",
				value:      "1",
				cooldown:   "PT15M",
			},
		},
	}

	return &profile
}

// Everything below here should be incorporated into TerraTest once they
// figure out their Azure support model they are working out with Microsoft

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

func GetMetricAlertsResource(t *testing.T, ruleName string) *insights.MetricAlertResource {
	metricAlertsResource, err := getMetricAlertsResourceE(ruleName)
	require.NoError(t, err)

	return metricAlertsResource
}

func getMetricAlertsResourceE(ruleName string) (*insights.MetricAlertResource, error) {
	client, err := getMetricAlertsClient()
	if err != nil {
		return nil, err
	}

	rule, err := client.Get(context.Background(), AzureResGroupName, ruleName)
	if err != nil {
		return nil, err
	}

	return &rule, nil
}

func getMetricAlertsClient() (*insights.MetricAlertsClient, error) {
	subID := os.Getenv(AzureSubscriptionID)

	metricAlertsClient := insights.NewMetricAlertsClient(subID)

	authorizer, err := azure.NewAuthorizer()
	if err != nil {
		return nil, err
	}

	metricAlertsClient.Authorizer = *authorizer

	return &metricAlertsClient, nil
}

func GetAutoscaleSettingsResource(t *testing.T, ruleName string) *insights.AutoscaleSettingResource {
	autoscaleSettingResource, err := getAutoscaleSettingsResourceE(ruleName)
	require.NoError(t, err)

	return autoscaleSettingResource
}

func getAutoscaleSettingsResourceE(ruleName string) (*insights.AutoscaleSettingResource, error) {
	client, err := getAutoscaleSettingsClient()
	if err != nil {
		return nil, err
	}

	autoscaleSetting, err := client.Get(context.Background(), AzureResGroupName, ruleName)
	if err != nil {
		return nil, err
	}

	return &autoscaleSetting, nil
}

func getAutoscaleSettingsClient() (*insights.AutoscaleSettingsClient, error) {
	subID := os.Getenv(AzureSubscriptionID)

	autoscaleSetting := insights.NewAutoscaleSettingsClient(subID)

	authorizer, err := azure.NewAuthorizer()
	if err != nil {
		return nil, err
	}

	autoscaleSetting.Authorizer = *authorizer

	return &autoscaleSetting, nil
}
