package monitoringstack

import (
	"testing"

	monv1 "github.com/rhobs/obo-prometheus-operator/pkg/apis/monitoring/v1"
	"github.com/rhobs/observability-operator/pkg/apis/monitoring/v1alpha1"
	"gotest.tools/v3/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestUpdateConditions(t *testing.T) {
	transitionTime := metav1.Now()

	tt := []struct {
		name                 string
		originalMSConditions []v1alpha1.Condition
		prometheus           monv1.Prometheus
		generation           int64
		recError             error
		expectedResults      []v1alpha1.Condition
		// flag to compare LastTransitionTime of each condition
		sameTransitionTimes bool
	}{
		{
			name:                 "empty conditions",
			originalMSConditions: []v1alpha1.Condition{},
			generation:           1,
			recError:             nil,
			prometheus:           monv1.Prometheus{},
			expectedResults: []v1alpha1.Condition{
				{
					Type:   v1alpha1.AvailableCondition,
					Status: v1alpha1.ConditionUnknown,
					Reason: "None",
				},
				{
					Type:   v1alpha1.ReconciledCondition,
					Status: v1alpha1.ConditionUnknown,
					Reason: "None",
				}},
		},
		{
			name: "conditions not changed when Prometheus Available",
			originalMSConditions: []v1alpha1.Condition{
				{
					Type:               v1alpha1.AvailableCondition,
					Status:             v1alpha1.ConditionTrue,
					ObservedGeneration: 1,
					Reason:             AvailableReason,
					Message:            AvailableMessage,
					LastTransitionTime: transitionTime,
				},
				{
					Type:               v1alpha1.ReconciledCondition,
					Status:             v1alpha1.ConditionTrue,
					ObservedGeneration: 1,
					Reason:             ReconciledReason,
					Message:            SuccessfullyReconciledMessage,
					LastTransitionTime: transitionTime,
				},
			},
			generation: 1,
			recError:   nil,
			prometheus: monv1.Prometheus{
				ObjectMeta: metav1.ObjectMeta{
					Generation: 1,
				},
				Status: monv1.PrometheusStatus{
					Conditions: []monv1.PrometheusCondition{
						{
							Type:               monv1.PrometheusAvailable,
							Status:             monv1.PrometheusConditionTrue,
							ObservedGeneration: 1,
						},
						{
							Type:               monv1.PrometheusReconciled,
							Status:             monv1.PrometheusConditionTrue,
							ObservedGeneration: 1,
						},
					}}},
			expectedResults: []v1alpha1.Condition{
				{
					Type:               v1alpha1.AvailableCondition,
					Status:             v1alpha1.ConditionTrue,
					ObservedGeneration: 1,
					Reason:             AvailableReason,
					Message:            AvailableMessage,
					LastTransitionTime: transitionTime,
				},
				{
					Type:               v1alpha1.ReconciledCondition,
					Status:             v1alpha1.ConditionTrue,
					ObservedGeneration: 1,
					Reason:             ReconciledReason,
					Message:            SuccessfullyReconciledMessage,
					LastTransitionTime: transitionTime,
				}},
			sameTransitionTimes: true,
		},
		{
			name: "cannot read Prometheus conditions",
			originalMSConditions: []v1alpha1.Condition{
				{
					Type:               v1alpha1.AvailableCondition,
					Status:             v1alpha1.ConditionTrue,
					ObservedGeneration: 1,
					Reason:             AvailableReason,
					Message:            AvailableMessage,
					LastTransitionTime: transitionTime,
				},
				{
					Type:               v1alpha1.ReconciledCondition,
					Status:             v1alpha1.ConditionTrue,
					ObservedGeneration: 1,
					Reason:             ReconciledReason,
					Message:            SuccessfullyReconciledMessage,
					LastTransitionTime: transitionTime,
				},
			},
			generation: 1,
			recError:   nil,
			prometheus: monv1.Prometheus{},
			expectedResults: []v1alpha1.Condition{
				{
					Type:               v1alpha1.AvailableCondition,
					Status:             v1alpha1.ConditionUnknown,
					ObservedGeneration: 1,
					Reason:             PrometheusNotAvailable,
					Message:            CannotReadPrometheusConditions,
				},
				{
					Type:               v1alpha1.ReconciledCondition,
					Status:             v1alpha1.ConditionUnknown,
					ObservedGeneration: 1,
					Reason:             PrometheusNotReconciled,
					Message:            CannotReadPrometheusConditions,
				}},
		},
		{
			name: "degraded Prometheus conditions",
			originalMSConditions: []v1alpha1.Condition{
				{
					Type:               v1alpha1.AvailableCondition,
					Status:             v1alpha1.ConditionTrue,
					ObservedGeneration: 1,
					Reason:             AvailableReason,
					Message:            AvailableMessage,
					LastTransitionTime: transitionTime,
				},
				{
					Type:               v1alpha1.ReconciledCondition,
					Status:             v1alpha1.ConditionTrue,
					ObservedGeneration: 1,
					Reason:             ReconciledReason,
					Message:            SuccessfullyReconciledMessage,
					LastTransitionTime: transitionTime,
				},
			},
			generation: 1,
			recError:   nil,
			prometheus: monv1.Prometheus{
				ObjectMeta: metav1.ObjectMeta{
					Generation: 1,
				},
				Status: monv1.PrometheusStatus{
					Conditions: []monv1.PrometheusCondition{
						{
							Type:               monv1.PrometheusAvailable,
							Status:             monv1.PrometheusConditionDegraded,
							ObservedGeneration: 1,
						},
						{
							Type:               monv1.PrometheusReconciled,
							Status:             monv1.PrometheusConditionDegraded,
							ObservedGeneration: 1,
						},
					}}},
			expectedResults: []v1alpha1.Condition{
				{
					Type:               v1alpha1.AvailableCondition,
					Status:             v1alpha1.ConditionFalse,
					ObservedGeneration: 1,
					Reason:             PrometheusDegraded,
				},
				{
					Type:               v1alpha1.ReconciledCondition,
					Status:             v1alpha1.ConditionFalse,
					ObservedGeneration: 1,
					Reason:             PrometheusNotReconciled,
				}},
		},
		{
			name: "Prometheus observed generation is different from the Prometheus generation",
			originalMSConditions: []v1alpha1.Condition{
				{
					Type:               v1alpha1.AvailableCondition,
					Status:             v1alpha1.ConditionTrue,
					ObservedGeneration: 2,
					Reason:             AvailableReason,
					Message:            AvailableMessage,
					LastTransitionTime: transitionTime,
				},
				{
					Type:               v1alpha1.ReconciledCondition,
					Status:             v1alpha1.ConditionTrue,
					ObservedGeneration: 2,
					Reason:             ReconciledReason,
					Message:            SuccessfullyReconciledMessage,
					LastTransitionTime: transitionTime,
				},
			},
			generation: 2,
			recError:   nil,
			prometheus: monv1.Prometheus{
				ObjectMeta: metav1.ObjectMeta{
					Generation: 3,
				},
				Status: monv1.PrometheusStatus{
					Conditions: []monv1.PrometheusCondition{
						{
							Type:               monv1.PrometheusAvailable,
							Status:             monv1.PrometheusConditionFalse,
							ObservedGeneration: 2,
						},
						{
							Type:               monv1.PrometheusReconciled,
							Status:             monv1.PrometheusConditionFalse,
							ObservedGeneration: 2,
						},
					}}},
			expectedResults: []v1alpha1.Condition{
				{
					Type:               v1alpha1.AvailableCondition,
					Status:             v1alpha1.ConditionTrue,
					ObservedGeneration: 2,
					Reason:             AvailableReason,
					Message:            AvailableMessage,
					LastTransitionTime: transitionTime,
				},
				{
					Type:               v1alpha1.ReconciledCondition,
					Status:             v1alpha1.ConditionTrue,
					ObservedGeneration: 2,
					Reason:             ReconciledReason,
					Message:            SuccessfullyReconciledMessage,
					LastTransitionTime: transitionTime,
				}},
			sameTransitionTimes: true,
		},
	}

	for _, test := range tt {
		res := updateConditions(test.originalMSConditions, test.prometheus, test.generation, test.recError)
		for _, c := range res {
			expectedC := getConditionByType(test.expectedResults, c.Type)
			assert.Check(t, expectedC.Equal(c), "%s - expected:\n %v\n and got:\n %v\n", test.name, expectedC, c)
			if test.sameTransitionTimes {
				assert.Check(t, expectedC.LastTransitionTime.Equal(&c.LastTransitionTime))
			} else {
				assert.Check(t, c.LastTransitionTime.After(transitionTime.Time))
			}
		}
	}
}

func getConditionByType(conditions []v1alpha1.Condition, ctype v1alpha1.ConditionType) *v1alpha1.Condition {
	for _, c := range conditions {
		if c.Type == ctype {
			return &c
		}
	}
	return nil
}
