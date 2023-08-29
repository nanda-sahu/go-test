// (C) Copyright 2021-2022 Hewlett Packard Enterprise Development LP

// Package terrors provides a set of base errors that can be used by
// applications to reduce boilerplate error code. All errors implement error
// wrapping and must be used as pointer type (as returned by their constructors).
//
// Whilst the base errors can be used directly within applications, this
// can lead to a lack of specificity in error handling that can cause problems.
// It is recommended that instead the base errors are embedded into error
// types that are specific to the application in question. For example:
//
//	 var VolumeNotFoundType = &VolumeNotFound{}
//
//	 type VolumeNotFound struct {
//	     *terrors.NotFound
//	 }
//
//	 func NewVolumeNotFound(id uuid.UUID, err error) *VolumeNotFound {
//	     return &VolumeNotFound{
//	         NotFound: terrors.NewNotFound("volume", "id", id.String()),
//		 }
//	 }
package terrors
