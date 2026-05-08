package metrics

import (
	"fmt"
	"io"
	"sort"
	"strings"
	"sync"
)

// 中文：Labels 定义当前包使用的数据结构或接口。
// English: Labels defines a data structure or interface used by this package.
type Labels map[string]string

// 中文：Registry 定义当前包使用的数据结构或接口。
// English: Registry defines a data structure or interface used by this package.
type Registry struct {
	// 中文：namespace 保存当前结构中的配置或数据值。
	// English: namespace stores a configuration or data value for this struct.
	namespace string
	// 中文：mu 保存当前结构中的配置或数据值。
	// English: mu stores a configuration or data value for this struct.
	mu sync.RWMutex
	// 中文：counters 保存当前结构中的配置或数据值。
	// English: counters stores a configuration or data value for this struct.
	counters map[string]float64
	// 中文：summaries 保存当前结构中的配置或数据值。
	// English: summaries stores a configuration or data value for this struct.
	summaries map[string]*summary
}

// 中文：summary 定义当前包使用的数据结构或接口。
// English: summary defines a data structure or interface used by this package.
type summary struct {
	// 中文：sum 保存当前结构中的配置或数据值。
	// English: sum stores a configuration or data value for this struct.
	sum float64
	// 中文：count 保存当前结构中的配置或数据值。
	// English: count stores a configuration or data value for this struct.
	count uint64
}

// 中文：CounterSnapshot 定义当前包使用的数据结构或接口。
// English: CounterSnapshot defines a data structure or interface used by this package.
type CounterSnapshot struct {
	// 中文：Name 保存当前结构中的配置或数据值。
	// English: Name stores a configuration or data value for this struct.
	Name string `json:"name"`
	// 中文：Labels 保存当前结构中的配置或数据值。
	// English: Labels stores a configuration or data value for this struct.
	Labels Labels `json:"labels,omitempty"`
	// 中文：Value 保存当前结构中的配置或数据值。
	// English: Value stores a configuration or data value for this struct.
	Value float64 `json:"value"`
}

// 中文：SummarySnapshot 定义当前包使用的数据结构或接口。
// English: SummarySnapshot defines a data structure or interface used by this package.
type SummarySnapshot struct {
	// 中文：Name 保存当前结构中的配置或数据值。
	// English: Name stores a configuration or data value for this struct.
	Name string `json:"name"`
	// 中文：Labels 保存当前结构中的配置或数据值。
	// English: Labels stores a configuration or data value for this struct.
	Labels Labels `json:"labels,omitempty"`
	// 中文：Sum 保存当前结构中的配置或数据值。
	// English: Sum stores a configuration or data value for this struct.
	Sum float64 `json:"sum"`
	// 中文：Count 保存当前结构中的配置或数据值。
	// English: Count stores a configuration or data value for this struct.
	Count uint64 `json:"count"`
	// 中文：Average 保存当前结构中的配置或数据值。
	// English: Average stores a configuration or data value for this struct.
	Average float64 `json:"average"`
}

// 中文：Snapshot 定义当前包使用的数据结构或接口。
// English: Snapshot defines a data structure or interface used by this package.
type Snapshot struct {
	// 中文：Namespace 保存当前结构中的配置或数据值。
	// English: Namespace stores a configuration or data value for this struct.
	Namespace string `json:"namespace"`
	// 中文：Counters 保存当前结构中的配置或数据值。
	// English: Counters stores a configuration or data value for this struct.
	Counters []CounterSnapshot `json:"counters"`
	// 中文：Summaries 保存当前结构中的配置或数据值。
	// English: Summaries stores a configuration or data value for this struct.
	Summaries []SummarySnapshot `json:"summaries"`
}

// 中文：NewRegistry 创建并返回对应组件实例。
// English: NewRegistry creates and returns the corresponding component instance.
func NewRegistry(namespace string) *Registry {
	namespace = strings.TrimSpace(namespace)
	if namespace == "" {
		namespace = "spiringo"
	}
	return &Registry{
		namespace: namespace,
		counters:  make(map[string]float64),
		summaries: make(map[string]*summary),
	}
}

// 中文：IncCounter 执行当前包中的对应流程。
// English: IncCounter executes the corresponding workflow in this package.
func (r *Registry) IncCounter(name string, labels Labels) {
	r.AddCounter(name, labels, 1)
}

// 中文：AddCounter 执行当前包中的对应流程。
// English: AddCounter executes the corresponding workflow in this package.
func (r *Registry) AddCounter(name string, labels Labels, delta float64) {
	if name == "" {
		return
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.counters[metricKey(name, labels)] += delta
}

// 中文：ObserveSummary 执行当前包中的对应流程。
// English: ObserveSummary executes the corresponding workflow in this package.
func (r *Registry) ObserveSummary(name string, labels Labels, value float64) {
	if name == "" {
		return
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	key := metricKey(name, labels)
	s, ok := r.summaries[key]
	if !ok {
		s = &summary{}
		r.summaries[key] = s
	}
	s.sum += value
	s.count++
}

// 中文：WritePrometheus 执行当前包中的对应流程。
// English: WritePrometheus executes the corresponding workflow in this package.
func (r *Registry) WritePrometheus(w io.Writer) error {
	r.mu.RLock()
	defer r.mu.RUnlock()

	counterNames := make(map[string]struct{})
	counterKeys := sortedKeys(r.counters)
	for _, key := range counterKeys {
		name, labels := splitMetricKey(key)
		fullName := r.fullName(name)
		if _, exists := counterNames[fullName]; !exists {
			if _, err := fmt.Fprintf(w, "# TYPE %s counter\n", fullName); err != nil {
				return err
			}
			counterNames[fullName] = struct{}{}
		}
		if _, err := fmt.Fprintf(w, "%s%s %g\n", fullName, formatLabels(labels), r.counters[key]); err != nil {
			return err
		}
	}

	summaryNames := make(map[string]struct{})
	summaryKeys := sortedKeys(r.summaries)
	for _, key := range summaryKeys {
		name, labels := splitMetricKey(key)
		fullName := r.fullName(name)
		if _, exists := summaryNames[fullName]; !exists {
			if _, err := fmt.Fprintf(w, "# TYPE %s summary\n", fullName); err != nil {
				return err
			}
			summaryNames[fullName] = struct{}{}
		}
		s := r.summaries[key]
		labelText := formatLabels(labels)
		if _, err := fmt.Fprintf(w, "%s_sum%s %g\n", fullName, labelText, s.sum); err != nil {
			return err
		}
		if _, err := fmt.Fprintf(w, "%s_count%s %d\n", fullName, labelText, s.count); err != nil {
			return err
		}
	}

	return nil
}

// 中文：Snapshot 执行当前包中的对应流程。
// English: Snapshot executes the corresponding workflow in this package.
func (r *Registry) Snapshot() Snapshot {
	r.mu.RLock()
	defer r.mu.RUnlock()

	snapshot := Snapshot{
		Namespace: r.namespace,
		Counters:  make([]CounterSnapshot, 0, len(r.counters)),
		Summaries: make([]SummarySnapshot, 0, len(r.summaries)),
	}
	for _, key := range sortedKeys(r.counters) {
		name, labels := splitMetricKey(key)
		snapshot.Counters = append(snapshot.Counters, CounterSnapshot{
			Name:   r.fullName(name),
			Labels: labels,
			Value:  r.counters[key],
		})
	}
	for _, key := range sortedKeys(r.summaries) {
		name, labels := splitMetricKey(key)
		s := r.summaries[key]
		avg := 0.0
		if s.count > 0 {
			avg = s.sum / float64(s.count)
		}
		snapshot.Summaries = append(snapshot.Summaries, SummarySnapshot{
			Name:    r.fullName(name),
			Labels:  labels,
			Sum:     s.sum,
			Count:   s.count,
			Average: avg,
		})
	}
	return snapshot
}

// 中文：WriteReport 执行当前包中的对应流程。
// English: WriteReport executes the corresponding workflow in this package.
func (r *Registry) WriteReport(w io.Writer) error {
	snapshot := r.Snapshot()
	if _, err := fmt.Fprintf(w, "# Metrics Report\n\nNamespace: `%s`\n\n", snapshot.Namespace); err != nil {
		return err
	}
	if _, err := fmt.Fprintln(w, "## Counters"); err != nil {
		return err
	}
	if len(snapshot.Counters) == 0 {
		if _, err := fmt.Fprintln(w, "\nNo counters recorded."); err != nil {
			return err
		}
	} else {
		if _, err := fmt.Fprintln(w, "\n| Name | Labels | Value |"); err != nil {
			return err
		}
		if _, err := fmt.Fprintln(w, "| --- | --- | ---: |"); err != nil {
			return err
		}
		for _, c := range snapshot.Counters {
			if _, err := fmt.Fprintf(w, "| `%s` | %s | %g |\n", c.Name, labelsForReport(c.Labels), c.Value); err != nil {
				return err
			}
		}
	}

	if _, err := fmt.Fprintln(w, "\n## Summaries"); err != nil {
		return err
	}
	if len(snapshot.Summaries) == 0 {
		_, err := fmt.Fprintln(w, "\nNo summaries recorded.")
		return err
	}
	if _, err := fmt.Fprintln(w, "\n| Name | Labels | Count | Sum | Average |"); err != nil {
		return err
	}
	if _, err := fmt.Fprintln(w, "| --- | --- | ---: | ---: | ---: |"); err != nil {
		return err
	}
	for _, s := range snapshot.Summaries {
		if _, err := fmt.Fprintf(w, "| `%s` | %s | %d | %.6f | %.6f |\n", s.Name, labelsForReport(s.Labels), s.Count, s.Sum, s.Average); err != nil {
			return err
		}
	}
	return nil
}

// 中文：fullName 执行当前包中的对应流程。
// English: fullName executes the corresponding workflow in this package.
func (r *Registry) fullName(name string) string {
	name = strings.TrimSpace(name)
	if r.namespace == "" {
		return name
	}
	return r.namespace + "_" + name
}

// 中文：metricKey 执行当前包中的对应流程。
// English: metricKey executes the corresponding workflow in this package.
func metricKey(name string, labels Labels) string {
	return name + "\xff" + encodeLabels(labels)
}

// 中文：splitMetricKey 执行当前包中的对应流程。
// English: splitMetricKey executes the corresponding workflow in this package.
func splitMetricKey(key string) (string, Labels) {
	parts := strings.SplitN(key, "\xff", 2)
	if len(parts) == 1 {
		return parts[0], nil
	}
	return parts[0], decodeLabels(parts[1])
}

// 中文：encodeLabels 执行当前包中的对应流程。
// English: encodeLabels executes the corresponding workflow in this package.
func encodeLabels(labels Labels) string {
	if len(labels) == 0 {
		return ""
	}
	keys := make([]string, 0, len(labels))
	for key := range labels {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	pairs := make([]string, 0, len(keys))
	for _, key := range keys {
		pairs = append(pairs, key+"="+labels[key])
	}
	return strings.Join(pairs, "\xff")
}

// 中文：decodeLabels 执行当前包中的对应流程。
// English: decodeLabels executes the corresponding workflow in this package.
func decodeLabels(encoded string) Labels {
	if encoded == "" {
		return nil
	}
	labels := Labels{}
	for _, pair := range strings.Split(encoded, "\xff") {
		key, value, ok := strings.Cut(pair, "=")
		if ok {
			labels[key] = value
		}
	}
	return labels
}

// 中文：formatLabels 执行当前包中的对应流程。
// English: formatLabels executes the corresponding workflow in this package.
func formatLabels(labels Labels) string {
	if len(labels) == 0 {
		return ""
	}
	keys := make([]string, 0, len(labels))
	for key := range labels {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	parts := make([]string, 0, len(keys))
	for _, key := range keys {
		parts = append(parts, fmt.Sprintf(`%s="%s"`, sanitizeLabelName(key), escapeLabelValue(labels[key])))
	}
	return "{" + strings.Join(parts, ",") + "}"
}

// 中文：sanitizeLabelName 执行当前包中的对应流程。
// English: sanitizeLabelName executes the corresponding workflow in this package.
func sanitizeLabelName(name string) string {
	var b strings.Builder
	for i, r := range name {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || r == '_' || (i > 0 && r >= '0' && r <= '9') {
			b.WriteRune(r)
			continue
		}
		b.WriteByte('_')
	}
	if b.Len() == 0 {
		return "label"
	}
	return b.String()
}

// 中文：escapeLabelValue 执行当前包中的对应流程。
// English: escapeLabelValue executes the corresponding workflow in this package.
func escapeLabelValue(value string) string {
	value = strings.ReplaceAll(value, `\`, `\\`)
	value = strings.ReplaceAll(value, "\n", `\n`)
	value = strings.ReplaceAll(value, `"`, `\"`)
	return value
}

// 中文：sortedKeys 执行当前包中的对应流程。
// English: sortedKeys executes the corresponding workflow in this package.
func sortedKeys[V any](m map[string]V) []string {
	keys := make([]string, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

// 中文：labelsForReport 执行当前包中的对应流程。
// English: labelsForReport executes the corresponding workflow in this package.
func labelsForReport(labels Labels) string {
	if len(labels) == 0 {
		return "-"
	}
	keys := make([]string, 0, len(labels))
	for key := range labels {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	parts := make([]string, 0, len(keys))
	for _, key := range keys {
		parts = append(parts, fmt.Sprintf("`%s=%s`", key, labels[key]))
	}
	return strings.Join(parts, "<br>")
}
