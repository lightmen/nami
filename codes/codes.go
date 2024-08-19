package codes

const (
	// OK is returned on success.
	OK int32 = 0

	// Canceled indicates the operation was canceled (typically by the caller).
	Canceled int32 = 1

	// Unknown error
	Unknown int32 = 2

	// InvalidArgument indicates an invalid argument.
	InvalidArgument int32 = 3

	// DeadlineExceeded means operation expired before completion.
	DeadlineExceeded int32 = 4

	// NotFound means some requested entity (e.g., file or directory) was
	// not found.
	NotFound int32 = 5

	// AlreadyExists means an attempt to create an entity failed because one
	// already exists.
	AlreadyExists int32 = 6

	// PermissionDenied indicates the caller does not have permission to
	// execute the specified operation.
	PermissionDenied int32 = 7

	// ResourceExhausted indicates some resource has been exhausted, perhaps
	// a per-user quota, or perhaps the entire file system is out of space.
	ResourceExhausted int32 = 8

	// FailedPrecondition indicates operation was rejected because the
	// system is not in a state required for the operation's execution.
	FailedPrecondition int32 = 9

	// Aborted indicates the operation was aborted, typically due to a
	// concurrency issue like sequencer check failures, transaction aborts,
	// etc.
	Aborted int32 = 10

	// OutOfRange means operation was attempted past the valid range.
	OutOfRange int32 = 11

	// Unimplemented indicates operation is not implemented or not
	// supported/enabled in this service.
	Unimplemented int32 = 12

	// Internal errors. Means some invarinami expected by underlying
	// system has been broken.
	Internal int32 = 13

	// Unavailable indicates the service is currently unavailable.
	Unavailable int32 = 14

	// DataLoss indicates unrecoverable data loss or corruption.
	DataLoss int32 = 15

	// Unauthenticated indicates the request does not have valid
	// authentication credentials for the operation.
	Unauthenticated int32 = 16

	// InvalidConfig indicates an invalid config.
	InvalidConfig int32 = 17

	//Authentication indicates an invalid Authentication.
	Authentication int32 = 18

	_maxint32 = 1000
)
