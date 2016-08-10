// +build linux
// +build amd64

package syscall

import (
	"internal"
	"unsafe"
)

const (
	sys_READ  = 0
	sys_WRITE = 1
	sys_OPEN  = 2
	sys_CLOSE = 3

	sys_MMAP = 9

	sys_BRK = 12

	sys_SOCKET  = 41
	sys_CONNECT = 42

	sys_BIND   = 49
	sys_LISTEN = 50

	sys_SETSOCKOPT = 54

	sys_CLONE = 56

	sys_EXIT = 60

	sys_UNLINK = 87

	sys_FUTEX = 202
)

const (
	O_ACCMODE   = 0x3
	O_APPEND    = 0x400
	O_ASYNC     = 0x2000
	O_CLOEXEC   = 0x80000
	O_CREAT     = 0x40
	O_DIRECT    = 0x4000
	O_DIRECTORY = 0x10000
	O_DSYNC     = 0x1000
	O_EXCL      = 0x80
	O_FSYNC     = 0x101000
	O_LARGEFILE = 0x0
	O_NDELAY    = 0x800
	O_NOATIME   = 0x40000
	O_NOCTTY    = 0x100
	O_NOFOLLOW  = 0x20000
	O_NONBLOCK  = 0x800
	O_RDONLY    = 0x0
	O_RDWR      = 0x2
	O_RSYNC     = 0x101000
	O_SYNC      = 0x101000
	O_TRUNC     = 0x200
	O_WRONLY    = 0x1
)

const (
	PROT_EXEC      = 0x4
	PROT_GROWSDOWN = 0x1000000
	PROT_GROWSUP   = 0x2000000
	PROT_NONE      = 0x0
	PROT_READ      = 0x1
	PROT_WRITE     = 0x2

	MAP_32BIT      = 0x40
	MAP_ANON       = 0x20
	MAP_ANONYMOUS  = 0x20
	MAP_DENYWRITE  = 0x800
	MAP_EXECUTABLE = 0x1000
	MAP_FILE       = 0x0
	MAP_FIXED      = 0x10
	MAP_GROWSDOWN  = 0x100
	MAP_HUGETLB    = 0x40000
	MAP_LOCKED     = 0x2000
	MAP_NONBLOCK   = 0x10000
	MAP_NORESERVE  = 0x4000
	MAP_POPULATE   = 0x8000
	MAP_PRIVATE    = 0x2
	MAP_SHARED     = 0x1
	MAP_STACK      = 0x20000
	MAP_TYPE       = 0xf
)

const (
	AF_ALG        = 0x26
	AF_APPLETALK  = 0x5
	AF_ASH        = 0x12
	AF_ATMPVC     = 0x8
	AF_ATMSVC     = 0x14
	AF_AX25       = 0x3
	AF_BLUETOOTH  = 0x1f
	AF_BRIDGE     = 0x7
	AF_CAIF       = 0x25
	AF_CAN        = 0x1d
	AF_DECnet     = 0xc
	AF_ECONET     = 0x13
	AF_FILE       = 0x1
	AF_IEEE802154 = 0x24
	AF_INET       = 0x2
	AF_INET6      = 0xa
	AF_IPX        = 0x4
	AF_IRDA       = 0x17
	AF_ISDN       = 0x22
	AF_IUCV       = 0x20
	AF_KEY        = 0xf
	AF_LLC        = 0x1a
	AF_LOCAL      = 0x1
	AF_MAX        = 0x27
	AF_NETBEUI    = 0xd
	AF_NETLINK    = 0x10
	AF_NETROM     = 0x6
	AF_PACKET     = 0x11
	AF_PHONET     = 0x23
	AF_PPPOX      = 0x18
	AF_RDS        = 0x15
	AF_ROSE       = 0xb
	AF_ROUTE      = 0x10
	AF_RXRPC      = 0x21
	AF_SECURITY   = 0xe
	AF_SNA        = 0x16
	AF_TIPC       = 0x1e
	AF_UNIX       = 0x1
	AF_UNSPEC     = 0x0
	AF_WANPIPE    = 0x19
	AF_X25        = 0x9
)

const (
	SOCK_CLOEXEC   = 0x80000
	SOCK_DCCP      = 0x6
	SOCK_DGRAM     = 0x2
	SOCK_NONBLOCK  = 0x800
	SOCK_PACKET    = 0xa
	SOCK_RAW       = 0x3
	SOCK_RDM       = 0x4
	SOCK_SEQPACKET = 0x5
	SOCK_STREAM    = 0x1
)

const (
	SOL_AAL    = 0x109
	SOL_ATM    = 0x108
	SOL_DECNET = 0x105
	SOL_ICMPV6 = 0x3a
	SOL_IP     = 0x0
	SOL_IPV6   = 0x29
	SOL_IRDA   = 0x10a
	SOL_PACKET = 0x107
	SOL_RAW    = 0xff
	SOL_SOCKET = 0x1
	SOL_TCP    = 0x6
	SOL_X25    = 0x106
)

const (
	SO_ACCEPTCONN                    = 0x1e
	SO_ATTACH_FILTER                 = 0x1a
	SO_BINDTODEVICE                  = 0x19
	SO_BROADCAST                     = 0x6
	SO_BSDCOMPAT                     = 0xe
	SO_DEBUG                         = 0x1
	SO_DETACH_FILTER                 = 0x1b
	SO_DOMAIN                        = 0x27
	SO_DONTROUTE                     = 0x5
	SO_ERROR                         = 0x4
	SO_KEEPALIVE                     = 0x9
	SO_LINGER                        = 0xd
	SO_MARK                          = 0x24
	SO_NO_CHECK                      = 0xb
	SO_OOBINLINE                     = 0xa
	SO_PASSCRED                      = 0x10
	SO_PASSSEC                       = 0x22
	SO_PEERCRED                      = 0x11
	SO_PEERNAME                      = 0x1c
	SO_PEERSEC                       = 0x1f
	SO_PRIORITY                      = 0xc
	SO_PROTOCOL                      = 0x26
	SO_RCVBUF                        = 0x8
	SO_RCVBUFFORCE                   = 0x21
	SO_RCVLOWAT                      = 0x12
	SO_RCVTIMEO                      = 0x14
	SO_REUSEADDR                     = 0x2
	SO_RXQ_OVFL                      = 0x28
	SO_SECURITY_AUTHENTICATION       = 0x16
	SO_SECURITY_ENCRYPTION_NETWORK   = 0x18
	SO_SECURITY_ENCRYPTION_TRANSPORT = 0x17
	SO_SNDBUF                        = 0x7
	SO_SNDBUFFORCE                   = 0x20
	SO_SNDLOWAT                      = 0x13
	SO_SNDTIMEO                      = 0x15
	SO_TIMESTAMP                     = 0x1d
	SO_TIMESTAMPING                  = 0x25
	SO_TIMESTAMPNS                   = 0x23
	SO_TYPE                          = 0x3
)

const (
	CLOCK_REALTIME  = 0
	CLOCK_MONOTONIC = 1
)

const (
	FUTEX_WAIT            = 0
	FUTEX_WAKE            = 1
	FUTEX_FD              = 2
	FUTEX_REQUEUE         = 3
	FUTEX_CMP_REQUEUE     = 4
	FUTEX_WAKE_OP         = 5
	FUTEX_LOCK_PI         = 6
	FUTEX_UNLOCK_PI       = 7
	FUTEX_TRYLOCK_PI      = 8
	FUTEX_WAIT_BITSET     = 9
	FUTEX_WAKE_BITSET     = 10
	FUTEX_WAIT_REQUEUE_PI = 11
	FUTEX_CMP_REQUEUE_PI  = 12

	FUTEX_PRIVATE_FLAG   = 128
	FUTEX_CLOCK_REALTIME = 256
)

// The following VDSO code is borrowed from:
// https://golang.org/src/runtime/vdso_linux_amd64.go
// https://git.kernel.org/cgit/linux/kernel/git/torvalds/linux.git/tree/Documentation/vDSO/parse_vdso.c

/* How to extract and insert information held in the st_info field.  */
func _ELF64_ST_BIND(val byte) byte { return val >> 4 }
func _ELF64_ST_TYPE(val byte) byte { return val & 0xf }

const (
	_AT_SYSINFO_EHDR = 33

	_PT_LOAD    = 1 /* Loadable program segment */
	_PT_DYNAMIC = 2 /* Dynamic linking information */

	_DT_NULL   = 0 /* Marks end of dynamic section */
	_DT_HASH   = 4 /* Dynamic symbol hash table */
	_DT_STRTAB = 5 /* Address of string table */
	_DT_SYMTAB = 6 /* Address of symbol table */
	_DT_VERSYM = 0x6ffffff0
	_DT_VERDEF = 0x6ffffffc

	_VER_FLG_BASE = 0x1 /* Version definition of file itself */

	_SHN_UNDEF = 0 /* Undefined section */

	_SHT_DYNSYM = 11 /* Dynamic linker symbol table */

	_STT_FUNC = 2 /* Symbol is a code object */

	_STB_GLOBAL = 1 /* Global symbol */
	_STB_WEAK   = 2 /* Weak symbol */

	_EI_NIDENT = 16
)

type elf64Sym struct {
	st_name  uint32
	st_info  byte
	st_other byte
	st_shndx uint16
	st_value uint64
	st_size  uint64
}

type elf64Verdef struct {
	vd_version uint16 /* Version revision */
	vd_flags   uint16 /* Version information */
	vd_ndx     uint16 /* Version Index */
	vd_cnt     uint16 /* Number of associated aux entries */
	vd_hash    uint32 /* Version name hash value */
	vd_aux     uint32 /* Offset in bytes to verdaux array */
	vd_next    uint32 /* Offset in bytes to next verdef entry */
}

type elf64Ehdr struct {
	e_ident     [_EI_NIDENT]byte /* Magic number and other info */
	e_type      uint16           /* Object file type */
	e_machine   uint16           /* Architecture */
	e_version   uint32           /* Object file version */
	e_entry     uint64           /* Entry point virtual address */
	e_phoff     uint64           /* Program header table file offset */
	e_shoff     uint64           /* Section header table file offset */
	e_flags     uint32           /* Processor-specific flags */
	e_ehsize    uint16           /* ELF header size in bytes */
	e_phentsize uint16           /* Program header table entry size */
	e_phnum     uint16           /* Program header table entry count */
	e_shentsize uint16           /* Section header table entry size */
	e_shnum     uint16           /* Section header table entry count */
	e_shstrndx  uint16           /* Section header string table index */
}

type elf64Phdr struct {
	p_type   uint32 /* Segment type */
	p_flags  uint32 /* Segment flags */
	p_offset uint64 /* Segment file offset */
	p_vaddr  uint64 /* Segment virtual address */
	p_paddr  uint64 /* Segment physical address */
	p_filesz uint64 /* Segment size in file */
	p_memsz  uint64 /* Segment size in memory */
	p_align  uint64 /* Segment alignment */
}

type elf64Dyn struct {
	d_tag int64  /* Dynamic entry type */
	d_val uint64 /* Integer value */
}

type elf64Verdaux struct {
	vda_name uint32 /* Version or dependency names */
	vda_next uint32 /* Offset in bytes to next verdaux entry */
}

type version_key struct {
	version  string
	ver_hash uint32
}

type vdso_info struct {
	valid bool

	/* Load information */
	load_addr   uintptr
	load_offset uintptr /* load_addr - recorded vaddr */

	/* Symbol table */
	symtab     *[1 << 32]elf64Sym
	symstrings *[1 << 32]byte
	chain      []uint32
	bucket     []uint32

	/* Version table */
	versym *[1 << 32]uint16
	verdef *elf64Verdef
}

type auxvEntry struct {
	typ uintptr // Type.
	val uintptr // Value.
}

type symbol_key struct {
	name     string
	sym_hash uint32
	ptr      unsafe.Pointer
}

var linux26 = version_key{"LINUX_2.6", 0x3ae75f6}

var (
	clock_gettime func(clkid int, tp *Timespec) uintptr
)

var sym_keys = []symbol_key{
	{"__vdso_clock_gettime", 0xd35ec75, unsafe.Pointer(&clock_gettime)},
}

func vdso_init_from_sysinfo_ehdr(info *vdso_info, hdr *elf64Ehdr) {
	info.valid = false
	info.load_addr = uintptr(unsafe.Pointer(hdr))

	p := unsafe.Pointer(info.load_addr + uintptr(hdr.e_phoff))

	// We need two things from the segment table: the load offset
	// and the dynamic table.
	var found_vaddr bool
	var dyn *[1 << 20]elf64Dyn
	for i := uint16(0); i < hdr.e_phnum; i++ {
		pt := (*elf64Phdr)(add(p, uintptr(i)*unsafe.Sizeof(elf64Phdr{})))
		switch pt.p_type {
		case _PT_LOAD:
			if !found_vaddr {
				found_vaddr = true
				info.load_offset = info.load_addr + uintptr(pt.p_offset-pt.p_vaddr)
			}
		case _PT_DYNAMIC:
			dyn = (*[1 << 20]elf64Dyn)(unsafe.Pointer(info.load_addr + uintptr(pt.p_offset)))
		}
	}

	if !found_vaddr || dyn == nil {
		return // Failed
	}

	// Fish out the useful bits of the dynamic table.

	var hash *[1 << 30]uint32
	hash = nil
	info.symstrings = nil
	info.symtab = nil
	info.versym = nil
	info.verdef = nil
	for i := 0; dyn[i].d_tag != _DT_NULL; i++ {
		dt := &dyn[i]
		p := info.load_offset + uintptr(dt.d_val)
		switch dt.d_tag {
		case _DT_STRTAB:
			info.symstrings = (*[1 << 32]byte)(unsafe.Pointer(p))
		case _DT_SYMTAB:
			info.symtab = (*[1 << 32]elf64Sym)(unsafe.Pointer(p))
		case _DT_HASH:
			hash = (*[1 << 30]uint32)(unsafe.Pointer(p))
		case _DT_VERSYM:
			info.versym = (*[1 << 32]uint16)(unsafe.Pointer(p))
		case _DT_VERDEF:
			info.verdef = (*elf64Verdef)(unsafe.Pointer(p))
		}
	}

	if info.symstrings == nil || info.symtab == nil || hash == nil {
		return // Failed
	}

	if info.verdef == nil {
		info.versym = nil
	}

	// Parse the hash table header.
	nbucket := hash[0]
	nchain := hash[1]
	info.bucket = hash[2 : 2+nbucket]
	info.chain = hash[2+nbucket : 2+nbucket+nchain]

	// That's all we need.
	info.valid = true
}

func vdso_find_version(info *vdso_info, ver *version_key) int32 {
	if !info.valid {
		return 0
	}
	def := info.verdef
	for {
		if def.vd_flags&_VER_FLG_BASE == 0 {
			aux := (*elf64Verdaux)(add(unsafe.Pointer(def), uintptr(def.vd_aux)))
			if def.vd_hash == ver.ver_hash && streq(ver.version, &info.symstrings[aux.vda_name]) {
				return int32(def.vd_ndx & 0x7fff)
			}
		}
		if def.vd_next == 0 {
			break
		}
		def = (*elf64Verdef)(add(unsafe.Pointer(def), uintptr(def.vd_next)))
	}
	return -1 // can not match any version
}

func add(p unsafe.Pointer, offset uintptr) unsafe.Pointer {
	return unsafe.Pointer(uintptr(p) + offset)
}

func streq(s string, cs *byte) bool {
	a := (*[1<<31 - 1]byte)(unsafe.Pointer(cs))
	for i := 0; i < len(s); i++ {
		g := s[i]
		c := a[i]
		if c == 0 || g != c {
			return false
		}
	}
	return a[len(s)] == 0
}

func vdso_parse_symbols(info *vdso_info, version int32) {
	if !info.valid {
		return
	}
	for _, k := range sym_keys {
		for chain := info.bucket[k.sym_hash%uint32(len(info.bucket))]; chain != 0; chain = info.chain[chain] {
			sym := &info.symtab[chain]
			typ := _ELF64_ST_TYPE(sym.st_info)
			bind := _ELF64_ST_BIND(sym.st_info)
			if typ != _STT_FUNC || bind != _STB_GLOBAL && bind != _STB_WEAK || sym.st_shndx == _SHN_UNDEF {
				continue
			}
			if !streq(k.name, &info.symstrings[sym.st_name]) {
				continue
			}
			// Check symbol version.
			if info.versym != nil && version != 0 && int32(info.versym[chain]&0x7fff) != version {
				continue
			}
			*(*uintptr)(k.ptr) = info.load_offset + uintptr(sym.st_value)
			break
		}
	}
}

func init() {
	auxv := (*[1<<31 - 1]auxvEntry)(unsafe.Pointer(internal.Auxv))
	for i := 0; ; i++ {
		e := &auxv[i]
		switch e.typ {
		case 0:
			panic("no valid _AT_SYSINFO_EHDR in AUXV")
		case _AT_SYSINFO_EHDR:
			if e.val == 0 {
				continue
			}
			var info vdso_info
			vdso_init_from_sysinfo_ehdr(&info, (*elf64Ehdr)(unsafe.Pointer(e.val)))
			ver := vdso_find_version(&info, &linux26)
			vdso_parse_symbols(&info, ver)
			for i := 0; i < len(sym_keys); i++ {
				if *(*uintptr)(sym_keys[i].ptr) == 0 {
					WriteString(2, sym_keys[i].name)
					WriteString(2, " = 0\n")
				}
			}
			return
		}
	}
}
