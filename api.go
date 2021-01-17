package webhdfs

const (
	PathPrefix = "/webhdfs/v1/"
)
const (
	OpOpen                          = "OPEN"
	OpGetFileStatus                 = "GETFILESTATUS"
	OpListStatus                    = "LISTSTATUS"
	OpListStatusBatch               = "LISTSTATUS_BATCH"
	OpGetContentSummary             = "GETCONTENTSUMMARY"
	OpGetQuotaUsage                 = "GETQUOTAUSAGE"
	OpGetFileChecksum               = "GETFILECHECKSUM"
	OpGetHomeDirectory              = "GETHOMEDIRECTORY"
	OpGetDelegationToken            = "GETDELEGATIONTOKEN"
	OpGetTrashRoot                  = "GETTRASHROOT"
	OpGetXAttr                      = "GETXATTRS"
	OpGetXAttrs                     = "GETXATTRS"
	OpGetAllXAttrs                  = "GETXATTRS"
	OpListXAttrs                    = "LISTXATTRS"
	OpCheckAccess                   = "CHECKACCESS"
	OpGetAllStoragePolicy           = "GETALLSTORAGEPOLICY"
	OpGetStoragePolicy              = "GETSTORAGEPOLICY"
	OpGetSnapshotDiff               = "GETSNAPSHOTDIFF"
	OpGetSnapshottableDirectoryList = "GETSNAPSHOTTABLEDIRECTORYLIST"
	OpGetFileBlockLocations         = "GETFILEBLOCKLOCATIONS"
	OpGetECPolicy                   = "GETECPOLICY"
	OpCreate                        = "CREATE"
	OpMkDirs                        = "MKDIRS"
	OpCreateSymlink                 = "CREATESYMLINK"
	OpRename                        = "RENAME"
	OpSetReplication                = "SETREPLICATION"
	OpSetOwner                      = "SETOWNER"
	OpSetPermission                 = "SETPERMISSION"
	OpSetTimes                      = "SETTIMES"
	OpRenewDelegationToken          = "RENEWDELEGATIONTOKEN"
	OpCancelDelegationToken         = "CANCELDELEGATIONTOKEN"
	OpAllowSnapshot                 = "ALLOWSNAPSHOT"
	OpDisallowSnapshot              = "DISALLOWSNAPSHOT"
	OpCreateSnapshot                = "CREATESNAPSHOT"
	OpRenameSnapshot                = "RENAMESNAPSHOT"
	OpSetXAttr                      = "SETXATTR"
	OpRemoveXAttr                   = "REMOVEXATTR"
	OpSetStoragePolicy              = "SETSTORAGEPOLICY"
	OpSatisfyStoragePolicy          = "SATISFYSTORAGEPOLICY"
	OpEnableECPolicy                = "ENABLEECPOLICY"
	OpDisableECPolicy               = "DISABLEECPOLICY"
	OpSetECPolicy                   = "SETECPOLICY"
	OpAppend                        = "APPEND"
	OpConcat                        = "CONCAT"
	OpTruncate                      = "TRUNCATE"
	OpUnsetStoragePolicy            = "UNSETSTORAGEPOLICY"
	OpUnsetECPolicy                 = "UNSETECPOLICY"
	OpDelete                        = "DELETE"
	OpDeleteSnapshot                = "DELETESNAPSHOT"
)
