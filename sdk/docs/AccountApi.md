# \AccountApi

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**Deposit**](AccountApi.md#Deposit) | **Get** /v1/deposit | Get Account Deposit Details



## Deposit

> AccountDepositDetails Deposit(ctx, account, assetCode)
Get Account Deposit Details

Gets the details for depositing currency into an account

### Required Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**account** | **string**|  | 
**assetCode** | **string**|  | 

### Return type

[**AccountDepositDetails**](AccountDepositDetails.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

