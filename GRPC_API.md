# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [api/community.proto](#api_community-proto)
    - [EmailIn](#-EmailIn)
    - [EmailOut](#-EmailOut)
    - [GetSchoolDataIn](#-GetSchoolDataIn)
    - [GetSchoolDataOut](#-GetSchoolDataOut)
    - [IsUserStaffOut](#-IsUserStaffOut)
    - [LoginIn](#-LoginIn)
    - [ParticipantChangeEvent](#-ParticipantChangeEvent)
    - [SearchPeer](#-SearchPeer)
    - [SearchPeersIn](#-SearchPeersIn)
    - [SearchPeersOut](#-SearchPeersOut)
    - [SendEduLinkingCodeIn](#-SendEduLinkingCodeIn)
  
    - [CommunityService](#-CommunityService)
  
- [Scalar Value Types](#scalar-value-types)



<a name="api_community-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## api/community.proto



<a name="-EmailIn"></a>

### EmailIn
Data for searching for matches in peers&#39; info


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| email | [string](#string) |  | User&#39;s E-mail address |






<a name="-EmailOut"></a>

### EmailOut
Response with found match


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| isExist | [bool](#bool) |  |  |






<a name="-GetSchoolDataIn"></a>

### GetSchoolDataIn



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| nickName | [string](#string) |  |  |






<a name="-GetSchoolDataOut"></a>

### GetSchoolDataOut



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| className | [string](#string) |  |  |
| parallelName | [string](#string) |  |  |






<a name="-IsUserStaffOut"></a>

### IsUserStaffOut



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| isStaff | [bool](#bool) |  |  |






<a name="-LoginIn"></a>

### LoginIn



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| login | [string](#string) |  |  |






<a name="-ParticipantChangeEvent"></a>

### ParticipantChangeEvent



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| login | [string](#string) |  |  |
| old_value_str | [string](#string) |  |  |
| old_value_int | [int32](#int32) |  |  |
| new_value_str | [string](#string) |  |  |
| new_value_int | [int32](#int32) |  |  |
| at | [google.protobuf.Timestamp](#google-protobuf-Timestamp) |  |  |






<a name="-SearchPeer"></a>

### SearchPeer



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| login | [string](#string) |  |  |






<a name="-SearchPeersIn"></a>

### SearchPeersIn



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| substring | [string](#string) |  |  |
| limit | [int64](#int64) |  |  |
| offset | [int64](#int64) |  |  |






<a name="-SearchPeersOut"></a>

### SearchPeersOut



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| searchPeers | [SearchPeer](#SearchPeer) | repeated |  |






<a name="-SendEduLinkingCodeIn"></a>

### SendEduLinkingCodeIn



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| login | [string](#string) |  |  |





 

 

 


<a name="-CommunityService"></a>

### CommunityService
Service with peers&#39; info from edu platform

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| IsPeerExist | [.EmailIn](#EmailIn) | [.EmailOut](#EmailOut) | Method for checking that the user is a school 21&#39;s student |
| SearchPeers | [.SearchPeersIn](#SearchPeersIn) | [.SearchPeersOut](#SearchPeersOut) |  |
| GetPeerSchoolData | [.GetSchoolDataIn](#GetSchoolDataIn) | [.GetSchoolDataOut](#GetSchoolDataOut) |  |
| isUserStaff | [.LoginIn](#LoginIn) | [.IsUserStaffOut](#IsUserStaffOut) |  |
| RunLoginsWorkerManually | [.google.protobuf.Empty](#google-protobuf-Empty) | [.google.protobuf.Empty](#google-protobuf-Empty) |  |
| SendEduLinkingCode | [.SendEduLinkingCodeIn](#SendEduLinkingCodeIn) | [.google.protobuf.Empty](#google-protobuf-Empty) |  |

 



## Scalar Value Types

| .proto Type | Notes | C++ | Java | Python | Go | C# | PHP | Ruby |
| ----------- | ----- | --- | ---- | ------ | -- | -- | --- | ---- |
| <a name="double" /> double |  | double | double | float | float64 | double | float | Float |
| <a name="float" /> float |  | float | float | float | float32 | float | float | Float |
| <a name="int32" /> int32 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint32 instead. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="int64" /> int64 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint64 instead. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="uint32" /> uint32 | Uses variable-length encoding. | uint32 | int | int/long | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="uint64" /> uint64 | Uses variable-length encoding. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum or Fixnum (as required) |
| <a name="sint32" /> sint32 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int32s. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sint64" /> sint64 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int64s. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="fixed32" /> fixed32 | Always four bytes. More efficient than uint32 if values are often greater than 2^28. | uint32 | int | int | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="fixed64" /> fixed64 | Always eight bytes. More efficient than uint64 if values are often greater than 2^56. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum |
| <a name="sfixed32" /> sfixed32 | Always four bytes. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sfixed64" /> sfixed64 | Always eight bytes. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="bool" /> bool |  | bool | boolean | boolean | bool | bool | boolean | TrueClass/FalseClass |
| <a name="string" /> string | A string must always contain UTF-8 encoded or 7-bit ASCII text. | string | String | str/unicode | string | string | string | String (UTF-8) |
| <a name="bytes" /> bytes | May contain any arbitrary sequence of bytes. | string | ByteString | str | []byte | ByteString | string | String (ASCII-8BIT) |

