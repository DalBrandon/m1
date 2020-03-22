// --------------------------------------------------------------------
// doc.go -- documentation for uuid package
//
// Created 2018-09-21 DLB
// --------------------------------------------------------------------

// The uuid package defines universally unique IDs that are suitable for database
// keys.  The primary type is a UUID, which it is a string of 32 hex
// characters without any puncuation marks or dashes.
//
// To create a new, completely and globally unique uuid in every way,
// call New(), as follows:
//
//    import "nsys/lib/uuid"
//    u1 := uuid.New()
//
// To get it's value, use String(), like so:
//
//    fmt.Printf("Value of this uuid: %s", u1.String())
//
// Once UUIDs are created, the are direcly compariable, for example:
//
//    u2 := uuid.New()   // Another, independent UUID
//    if u1 == u2 {
//        fmt.Printf("The world has ended!")  // This should NEVER happen.
//    }
//
// An all-zero UUID is special.  It is used to indicate an uninitialized or unspecified
// UUID.  Special care is taken with uninitialized UUIDs to make sure they
// behave as all-zero uuids.  That is, an all-zero uuid and an un-initialized
// uuid are equal, and you can do the following:
//
//    var u3 uuid.UUID      // Declare an unitialized UUID
//    u4 := uuid.ForceStr("00000000000000000000000000000000")  // Create a zero uuid explicitily
//    u5 := uuid.Zero()    // Create a zero uuid with the built in function
//    if u3 != u4 || u4 != u5 || u3 != u5 {
//        fmt.Printf("This should never happen because u3, u4 and u5 are all equal to zero.")
//    }
//
// These functions help in converting a string into a uuid:
//
//    uuid.FromString()  -- Returns error if string is not exactly 32 hex chars
//    uuid.FromString0() -- Like FromString() but allows blank input
//    uuid.ForceStr()    -- Does not return an error, but returns zero uuid on bad input
//    uuid.IsUuidString() -- Returns true if the input string is exactly 32 hex chars
//
// There functions are provided for dealing with zero uuids:
//
//    uuid.Zero() -- returns a zero uuid
//    uuid.IsZero() -- returns true if the uuid is zero
//
// In addition, UUID satisfies the Json and Gob marshalling interfaces -- so that UUID
// can be used in data structures without worry about encoding issues.  The Json decoding
// is very liberal about what it will accecpt: it is case insensitive, and it deals with
// blanks without problem.  Likewise the sql db driver interface is provided so that uuids
// can work directly with sql databases.
//
// Finally, note that uuids support the Stringer interface, so it is okay to do the
// following:
//
//    u := uuid.New()
//    text := fmt.Sprintf("The value of the uuid is %s.", u)
//
package uuid
