# \AuthorizationApi

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**Authenticate**](AuthorizationApi.md#Authenticate) | **Post** /v1/authorizations | Authenticate A Challenge
[**RequestAChallenge**](AuthorizationApi.md#RequestAChallenge) | **Get** /v1/authorizations | Request A Challenge



## Authenticate

> AuthToken Authenticate(ctx, account, challengeTransaction)
Authenticate A Challenge

Authenticates the user with a signed JWT token

### Required Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**account** | **string**|  | 
**challengeTransaction** | [**ChallengeTransaction**](ChallengeTransaction.md)|  | 

### Return type

[**AuthToken**](AuthToken.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## RequestAChallenge

> ChallengeTransaction RequestAChallenge(ctx, account)
Request A Challenge

Requests a challenge transaction for the client to sign. Once signed, the challenge transaction is to be resubmitted back to the api.

### Required Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**account** | **string**|  | 

### Return type

[**ChallengeTransaction**](ChallengeTransaction.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

