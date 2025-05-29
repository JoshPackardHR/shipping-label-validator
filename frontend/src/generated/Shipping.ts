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

import { ShippingError, ShippingRequest, ShippingResponse } from "./data-contracts";
import { ContentType, HttpClient, RequestParams } from "./http-client";

export class Shipping<SecurityDataType = unknown> extends HttpClient<SecurityDataType> {
  /**
   * @description check a shipping label
   *
   * @tags shipping
   * @name LabelCheckCreate
   * @summary Check a shipping label
   * @request POST:/shipping/label/check
   */
  labelCheckCreate = (requestBody: ShippingRequest, params: RequestParams = {}) =>
    this.request<ShippingResponse, ShippingError>({
      path: `/shipping/label/check`,
      method: "POST",
      body: requestBody,
      type: ContentType.Json,
      format: "json",
      ...params,
    });
}
