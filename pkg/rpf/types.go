package rpf

type RPFHeader struct {
	Version     int32
	TOCSize     int32
	EntryCount  int32
	Unknown     int32
	IsEncrypted int32
}

type DirectoryEntry struct {
	NameOffset        int32
	Flags             int32
	ContentEntryIndex uint32
	ContentEntryCount uint32
}

type FileEntry struct {
	NameOffset   int32
	Size         uint32 // Changed from int32 to uint32
	Offset       [3]byte
	ResourceType byte
	Flags        uint32
}
