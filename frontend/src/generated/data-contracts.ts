/* eslint-disable */
/* tslint:disable */
// @ts-nocheck
/*
 * ---------------------------------------------------------------
 * ## THIS FILE WAS GENERATED VIA SWAGGER-TYPESCRIPT-API        ##
 * ##                                                           ##
 * ## AUTHOR: acacode                                           ##
 * ## SOURCE: https://github.com/acacode/swagger-typescript-api ##
 * ---------------------------------------------------------------
 */

export interface UpsAddress {
  addressLine1?: string;
  addressLine2?: string;
  city?: string;
  country?: string;
  countryCode?: string;
  postalCode?: string;
  stateProvince?: string;
}

export interface UpsPackageAddress {
  address?: UpsAddress;
  attentionName?: string;
  name?: string;
  type?: string;
}

export interface ValidationError {
  error?: string;
}

export interface ValidationRequest {
  image?: string;
  trackingNumber?: string;
}

export interface ValidationResponse {
  result?: ValidationResult;
}

export interface ValidationResult {
  expectedAddress?: UpsPackageAddress;
  scannedAddress?: UpsAddress;
  valid?: boolean;
}
