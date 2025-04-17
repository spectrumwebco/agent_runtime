# Kata Containers configuration
[hypervisor.qemu]
path = "/usr/bin/qemu-system-x86_64"
kernel = "/usr/share/kata-containers/vmlinux.container"
initrd = "/usr/share/kata-containers/kata-containers-initrd.img"
machine_type = "q35"
firmware = ""
default_vcpus = 1
default_memory = 1024
default_maxvcpus = 0
default_bridges = 1
block_device_driver = "virtio-scsi"
disable_block_device_use = false
enable_iothreads = false
enable_mem_prealloc = false
enable_hugepages = false
enable_swap = false
enable_debug = false
disable_nesting_checks = false
msize_9p = 8192
use_vsock = true
hotplug_vfio_on_root_bus = false
disable_vhost_net = true
entropy_source = "/dev/urandom"
enable_template = false

[runtime]
enable_cpu_quota = true
internetworking_model = "tcfilter"
disable_new_netns = false
sandbox_cgroup_only = false
experimental = ${enable_runtime_rs ? "[\"runtime_rs\"]" : "[]"}

[agent.kata]
enable_tracing = false
enable_debug = false
dial_timeout = 0
debug_console_enabled = false
debug_console_vport = 0
enable_mem_hotplug = false
kernel_modules = []
container_pipe_size = 0
use_vsock = true

[netmon]
enable_netmon = false
netmon_path = "/usr/libexec/kata-containers/kata-netmon"
netmon_socket_path = "/run/vc/sockets/netmon.sock"
netmon_container_pipe_size = 0

[image]
service_offload = false
service_image = ""

[shim.kata]
enable_tracing = false
tracing_socket_path = ""
disable_guest_seccomp = false

[mem_agent]
enable_mem_agent = ${enable_mem_agent}
