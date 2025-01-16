package ebpfGenerated

//go:generate go run github.com/cilium/ebpf/cmd/bpf2go counter ../ebpf_module_kernel/counter.ebpf.c --pkg counterebpf
//go:generate go run github.com/cilium/ebpf/cmd/bpf2go filterip ../ebpf_module_kernel/filter_ip.ebpf.c --pkg filterip
//go:generate go run github.com/cilium/ebpf/cmd/bpf2go filtermac ../ebpf_module_kernel/filter_mac.ebpf.c --pkg filtermac
//go:generate go run github.com/cilium/ebpf/cmd/bpf2go dispacher ../ebpf_module_kernel/dispacher.ebpf.c --pkg dispacher
