package component

import (
	"fmt"

	"github.com/3scale/3scale-operator/pkg/assets"
	"github.com/3scale/3scale-operator/pkg/common"
	monitoringv1 "github.com/coreos/prometheus-operator/pkg/apis/monitoring/v1"
	grafanav1alpha1 "github.com/integr8ly/grafana-operator/v3/pkg/apis/integreatly/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func (zync *Zync) ZyncPodMonitor() *monitoringv1.PodMonitor {
	return &monitoringv1.PodMonitor{
		ObjectMeta: metav1.ObjectMeta{
			Name:   "zync",
			Labels: zync.Options.CommonZyncLabels,
		},
		Spec: monitoringv1.PodMonitorSpec{
			PodMetricsEndpoints: []monitoringv1.PodMetricsEndpoint{{
				Port:   "metrics",
				Path:   "/metrics",
				Scheme: "http",
			}},
			Selector: metav1.LabelSelector{
				MatchLabels: zync.Options.CommonZyncLabels,
			},
		},
	}
}

func (zync *Zync) ZyncQuePodMonitor() *monitoringv1.PodMonitor {
	return &monitoringv1.PodMonitor{
		ObjectMeta: metav1.ObjectMeta{
			Name:   "zync-que",
			Labels: zync.Options.CommonZyncQueLabels,
		},
		Spec: monitoringv1.PodMonitorSpec{
			PodMetricsEndpoints: []monitoringv1.PodMetricsEndpoint{{
				Port:   "metrics",
				Path:   "/metrics",
				Scheme: "http",
			}},
			Selector: metav1.LabelSelector{
				MatchLabels: zync.Options.CommonZyncQueLabels,
			},
		},
	}
}

func ZyncGrafanaDashboard(ns string) *grafanav1alpha1.GrafanaDashboard {
	data := &struct {
		Namespace string
	}{
		ns,
	}
	return &grafanav1alpha1.GrafanaDashboard{
		ObjectMeta: metav1.ObjectMeta{
			Name: "zync",
			Labels: map[string]string{
				"monitoring-key": common.MonitoringKey,
			},
		},
		Spec: grafanav1alpha1.GrafanaDashboardSpec{
			Json: assets.TemplateAsset("monitoring/zync-grafana-dashboard-1.json.tpl", data),
			Name: fmt.Sprintf("%s/zync-grafana-dashboard-1.json", ns),
		},
	}
}

func ZyncPrometheusRules(ns string) *monitoringv1.PrometheusRule {
	return &monitoringv1.PrometheusRule{
		ObjectMeta: metav1.ObjectMeta{
			Name: "zync",
			Labels: map[string]string{
				"prometheus": "application-monitoring",
				"role":       "alert-rules",
			},
		},
		Spec: monitoringv1.PrometheusRuleSpec{
			Groups: []monitoringv1.RuleGroup{
				{
					Name: fmt.Sprintf("%s/zync.rules", ns),
					Rules: []monitoringv1.Rule{
						{
							Alert: "PumaWorkersRunningLow",
							Annotations: map[string]string{
								"summary":     "{{$labels.container_name}} replica controller on {{$labels.namespace}}: Has less than 5 puma workers in the last 5 minutes",
								"description": "{{$labels.container_name}} replica controller on {{$labels.namespace}} project: Has less than 5 puma workers in the last 5 minutes",
							},
							Expr: intstr.FromString(fmt.Sprintf(`avg_over_time(puma_running{job="zync-monitoring",namespace="%s"} [5m]) < 5`, ns)),
							For:  "30m",
							Labels: map[string]string{
								"severity": "critical",
							},
						},
					},
				},
			},
		},
	}
}

func ZyncQuePrometheusRules(ns string) *monitoringv1.PrometheusRule {
	return &monitoringv1.PrometheusRule{
		ObjectMeta: metav1.ObjectMeta{
			Name: "zync-que",
			Labels: map[string]string{
				"prometheus": "application-monitoring",
				"role":       "alert-rules",
			},
		},
		Spec: monitoringv1.PrometheusRuleSpec{
			Groups: []monitoringv1.RuleGroup{
				{
					Name: fmt.Sprintf("%s/zync-que.rules", ns),
					Rules: []monitoringv1.Rule{
						{
							Alert: "QueWorkersRunningLow",
							Annotations: map[string]string{
								"summary":     "{{$labels.container_name}} replica controller on {{$labels.namespace}}: Has less than 5 que workers in the last 5 minutes",
								"description": "{{$labels.container_name}} replica controller on {{$labels.namespace}} project: Has less than 5 que workers in the last 5 minutes",
							},
							Expr: intstr.FromString(fmt.Sprintf(`avg_over_time(que_workers_total{job="zync-que-monitoring",namespace="%s"} [5m]) < 5`, ns)),
							For:  "30m",
							Labels: map[string]string{
								"severity": "critical",
							},
						},
					},
				},
			},
		},
	}
}