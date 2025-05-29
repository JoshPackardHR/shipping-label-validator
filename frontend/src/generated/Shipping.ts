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

import { ValidationError, ValidationRequest, ValidationResponse } from "./data-contracts";
import { ContentType, HttpClient, RequestParams } from "./http-client";

export class Shipping<SecurityDataType = unknown> extends HttpClient<SecurityDataType> {
  /**
   * @description check a shipping label
   *
   * @tags shipping
   * @name LabelValidateCreate
   * @summary Check a shipping label
   * @request POST:/shipping/label/validate
   */
  labelValidateCreate = (requestBody: ValidationRequest, params: RequestParams = {}) =>
    this.request<ValidationResponse, ValidationError>({
      path: `/shipping/label/validate`,
      method: "POST",
      body: requestBody,
      type: ContentType.Json,
      format: "json",
      ...params,
    });
}
