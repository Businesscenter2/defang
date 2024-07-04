// protos/v1/fabric.proto

// @generated by protoc-gen-es v1.10.0
// @generated from file io/defang/v1/fabric.proto (package io.defang.v1, syntax proto3)
/* eslint-disable */
// @ts-nocheck

import { proto3, Struct, Timestamp } from "@bufbuild/protobuf";

/**
 * @generated from enum io.defang.v1.Platform
 */
export const Platform = /*@__PURE__*/ proto3.makeEnum(
  "io.defang.v1.Platform",
  [
    {no: 0, name: "LINUX_AMD64"},
    {no: 1, name: "LINUX_ARM64"},
    {no: 2, name: "LINUX_ANY"},
  ],
);

/**
 * @generated from enum io.defang.v1.Protocol
 */
export const Protocol = /*@__PURE__*/ proto3.makeEnum(
  "io.defang.v1.Protocol",
  [
    {no: 0, name: "ANY"},
    {no: 1, name: "UDP"},
    {no: 2, name: "TCP"},
    {no: 3, name: "HTTP"},
    {no: 4, name: "HTTP2"},
    {no: 5, name: "GRPC"},
  ],
);

/**
 * @generated from enum io.defang.v1.Mode
 */
export const Mode = /*@__PURE__*/ proto3.makeEnum(
  "io.defang.v1.Mode",
  [
    {no: 0, name: "HOST"},
    {no: 1, name: "INGRESS"},
  ],
);

/**
 * @generated from enum io.defang.v1.Network
 */
export const Network = /*@__PURE__*/ proto3.makeEnum(
  "io.defang.v1.Network",
  [
    {no: 0, name: "UNSPECIFIED"},
    {no: 1, name: "PRIVATE"},
    {no: 2, name: "PUBLIC"},
  ],
);

/**
 * @generated from message io.defang.v1.RestartRequest
 */
export const RestartRequest = /*@__PURE__*/ proto3.makeMessageType(
  "io.defang.v1.RestartRequest",
  () => [
    { no: 1, name: "services", kind: "scalar", T: 9 /* ScalarType.STRING */, repeated: true },
    { no: 2, name: "project", kind: "scalar", T: 9 /* ScalarType.STRING */ },
  ],
);

/**
 * @generated from message io.defang.v1.TrackRequest
 */
export const TrackRequest = /*@__PURE__*/ proto3.makeMessageType(
  "io.defang.v1.TrackRequest",
  () => [
    { no: 1, name: "anon_id", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 2, name: "event", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 3, name: "properties", kind: "map", K: 9 /* ScalarType.STRING */, V: {kind: "scalar", T: 9 /* ScalarType.STRING */} },
    { no: 4, name: "os", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 5, name: "arch", kind: "scalar", T: 9 /* ScalarType.STRING */ },
  ],
);

/**
 * @generated from message io.defang.v1.DeployRequest
 */
export const DeployRequest = /*@__PURE__*/ proto3.makeMessageType(
  "io.defang.v1.DeployRequest",
  () => [
    { no: 1, name: "services", kind: "message", T: Service, repeated: true },
    { no: 2, name: "project", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 3, name: "compose", kind: "message", T: Struct },
  ],
);

/**
 * @generated from message io.defang.v1.DeployResponse
 */
export const DeployResponse = /*@__PURE__*/ proto3.makeMessageType(
  "io.defang.v1.DeployResponse",
  () => [
    { no: 1, name: "services", kind: "message", T: ServiceInfo, repeated: true },
    { no: 2, name: "etag", kind: "scalar", T: 9 /* ScalarType.STRING */ },
  ],
);

/**
 * @generated from message io.defang.v1.DeleteRequest
 */
export const DeleteRequest = /*@__PURE__*/ proto3.makeMessageType(
  "io.defang.v1.DeleteRequest",
  () => [
    { no: 1, name: "names", kind: "scalar", T: 9 /* ScalarType.STRING */, repeated: true },
  ],
);

/**
 * @generated from message io.defang.v1.DeleteResponse
 */
export const DeleteResponse = /*@__PURE__*/ proto3.makeMessageType(
  "io.defang.v1.DeleteResponse",
  () => [
    { no: 1, name: "etag", kind: "scalar", T: 9 /* ScalarType.STRING */ },
  ],
);

/**
 * @generated from message io.defang.v1.GenerateFilesRequest
 */
export const GenerateFilesRequest = /*@__PURE__*/ proto3.makeMessageType(
  "io.defang.v1.GenerateFilesRequest",
  () => [
    { no: 1, name: "prompt", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 2, name: "language", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 3, name: "agree_tos", kind: "scalar", T: 8 /* ScalarType.BOOL */ },
  ],
);

/**
 * @generated from message io.defang.v1.File
 */
export const File = /*@__PURE__*/ proto3.makeMessageType(
  "io.defang.v1.File",
  () => [
    { no: 1, name: "name", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 2, name: "content", kind: "scalar", T: 9 /* ScalarType.STRING */ },
  ],
);

/**
 * @generated from message io.defang.v1.GenerateFilesResponse
 */
export const GenerateFilesResponse = /*@__PURE__*/ proto3.makeMessageType(
  "io.defang.v1.GenerateFilesResponse",
  () => [
    { no: 1, name: "files", kind: "message", T: File, repeated: true },
  ],
);

/**
 * @generated from message io.defang.v1.StartGenerateResponse
 */
export const StartGenerateResponse = /*@__PURE__*/ proto3.makeMessageType(
  "io.defang.v1.StartGenerateResponse",
  () => [
    { no: 1, name: "uuid", kind: "scalar", T: 9 /* ScalarType.STRING */ },
  ],
);

/**
 * @generated from message io.defang.v1.GenerateStatusRequest
 */
export const GenerateStatusRequest = /*@__PURE__*/ proto3.makeMessageType(
  "io.defang.v1.GenerateStatusRequest",
  () => [
    { no: 1, name: "uuid", kind: "scalar", T: 9 /* ScalarType.STRING */ },
  ],
);

/**
 * @generated from message io.defang.v1.UploadURLRequest
 */
export const UploadURLRequest = /*@__PURE__*/ proto3.makeMessageType(
  "io.defang.v1.UploadURLRequest",
  () => [
    { no: 1, name: "digest", kind: "scalar", T: 9 /* ScalarType.STRING */ },
  ],
);

/**
 * @generated from message io.defang.v1.UploadURLResponse
 */
export const UploadURLResponse = /*@__PURE__*/ proto3.makeMessageType(
  "io.defang.v1.UploadURLResponse",
  () => [
    { no: 1, name: "url", kind: "scalar", T: 9 /* ScalarType.STRING */ },
  ],
);

/**
 * @generated from message io.defang.v1.ServiceInfo
 */
export const ServiceInfo = /*@__PURE__*/ proto3.makeMessageType(
  "io.defang.v1.ServiceInfo",
  () => [
    { no: 1, name: "service", kind: "message", T: ServiceID },
    { no: 2, name: "endpoints", kind: "scalar", T: 9 /* ScalarType.STRING */, repeated: true },
    { no: 3, name: "project", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 4, name: "etag", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 5, name: "status", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 6, name: "nat_ips", kind: "scalar", T: 9 /* ScalarType.STRING */, repeated: true },
    { no: 7, name: "lb_ips", kind: "scalar", T: 9 /* ScalarType.STRING */, repeated: true },
    { no: 8, name: "private_fqdn", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 9, name: "public_fqdn", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 10, name: "created_at", kind: "message", T: Timestamp },
    { no: 11, name: "updated_at", kind: "message", T: Timestamp },
    { no: 12, name: "zone_id", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 13, name: "use_acme_cert", kind: "scalar", T: 8 /* ScalarType.BOOL */ },
    { no: 15, name: "domainname", kind: "scalar", T: 9 /* ScalarType.STRING */ },
  ],
);

/**
 * @generated from message io.defang.v1.Secrets
 */
export const Secrets = /*@__PURE__*/ proto3.makeMessageType(
  "io.defang.v1.Secrets",
  () => [
    { no: 1, name: "names", kind: "scalar", T: 9 /* ScalarType.STRING */, repeated: true },
    { no: 2, name: "project", kind: "scalar", T: 9 /* ScalarType.STRING */ },
  ],
);

/**
 * @generated from message io.defang.v1.SecretValue
 */
export const SecretValue = /*@__PURE__*/ proto3.makeMessageType(
  "io.defang.v1.SecretValue",
  () => [
    { no: 1, name: "name", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 2, name: "value", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 3, name: "project", kind: "scalar", T: 9 /* ScalarType.STRING */ },
  ],
);

/**
 * @generated from message io.defang.v1.TokenRequest
 */
export const TokenRequest = /*@__PURE__*/ proto3.makeMessageType(
  "io.defang.v1.TokenRequest",
  () => [
    { no: 1, name: "tenant", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 2, name: "auth_code", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 3, name: "scope", kind: "scalar", T: 9 /* ScalarType.STRING */, repeated: true },
    { no: 4, name: "assertion", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 5, name: "expires_in", kind: "scalar", T: 13 /* ScalarType.UINT32 */ },
    { no: 6, name: "anon_id", kind: "scalar", T: 9 /* ScalarType.STRING */ },
  ],
);

/**
 * @generated from message io.defang.v1.TokenResponse
 */
export const TokenResponse = /*@__PURE__*/ proto3.makeMessageType(
  "io.defang.v1.TokenResponse",
  () => [
    { no: 1, name: "access_token", kind: "scalar", T: 9 /* ScalarType.STRING */ },
  ],
);

/**
 * @generated from message io.defang.v1.Status
 */
export const Status = /*@__PURE__*/ proto3.makeMessageType(
  "io.defang.v1.Status",
  () => [
    { no: 1, name: "version", kind: "scalar", T: 9 /* ScalarType.STRING */ },
  ],
);

/**
 * @generated from message io.defang.v1.Version
 */
export const Version = /*@__PURE__*/ proto3.makeMessageType(
  "io.defang.v1.Version",
  () => [
    { no: 1, name: "fabric", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 3, name: "cli_min", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 4, name: "pulumi_min", kind: "scalar", T: 9 /* ScalarType.STRING */ },
  ],
);

/**
 * @generated from message io.defang.v1.TailRequest
 */
export const TailRequest = /*@__PURE__*/ proto3.makeMessageType(
  "io.defang.v1.TailRequest",
  () => [
    { no: 1, name: "services", kind: "scalar", T: 9 /* ScalarType.STRING */, repeated: true },
    { no: 2, name: "since", kind: "message", T: Timestamp },
    { no: 3, name: "etag", kind: "scalar", T: 9 /* ScalarType.STRING */ },
  ],
);

/**
 * @generated from message io.defang.v1.LogEntry
 */
export const LogEntry = /*@__PURE__*/ proto3.makeMessageType(
  "io.defang.v1.LogEntry",
  () => [
    { no: 1, name: "message", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 2, name: "timestamp", kind: "message", T: Timestamp },
    { no: 3, name: "stderr", kind: "scalar", T: 8 /* ScalarType.BOOL */ },
    { no: 4, name: "service", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 5, name: "etag", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 6, name: "host", kind: "scalar", T: 9 /* ScalarType.STRING */ },
  ],
);

/**
 * @generated from message io.defang.v1.TailResponse
 */
export const TailResponse = /*@__PURE__*/ proto3.makeMessageType(
  "io.defang.v1.TailResponse",
  () => [
    { no: 2, name: "entries", kind: "message", T: LogEntry, repeated: true },
    { no: 3, name: "service", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 4, name: "etag", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 5, name: "host", kind: "scalar", T: 9 /* ScalarType.STRING */ },
  ],
);

/**
 * @generated from message io.defang.v1.ListServicesResponse
 */
export const ListServicesResponse = /*@__PURE__*/ proto3.makeMessageType(
  "io.defang.v1.ListServicesResponse",
  () => [
    { no: 1, name: "services", kind: "message", T: ServiceInfo, repeated: true },
    { no: 2, name: "project", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 3, name: "compose", kind: "message", T: Struct },
  ],
);

/**
 * TODO: internal message; move to a separate proto file
 *
 * @generated from message io.defang.v1.ProjectUpdate
 */
export const ProjectUpdate = /*@__PURE__*/ proto3.makeMessageType(
  "io.defang.v1.ProjectUpdate",
  () => [
    { no: 1, name: "services", kind: "message", T: ServiceInfo, repeated: true },
    { no: 2, name: "alb_arn", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 3, name: "compose", kind: "message", T: Struct },
  ],
);

/**
 * @generated from message io.defang.v1.ServiceID
 */
export const ServiceID = /*@__PURE__*/ proto3.makeMessageType(
  "io.defang.v1.ServiceID",
  () => [
    { no: 1, name: "name", kind: "scalar", T: 9 /* ScalarType.STRING */ },
  ],
);

/**
 * @generated from message io.defang.v1.Device
 */
export const Device = /*@__PURE__*/ proto3.makeMessageType(
  "io.defang.v1.Device",
  () => [
    { no: 1, name: "capabilities", kind: "scalar", T: 9 /* ScalarType.STRING */, repeated: true },
    { no: 2, name: "driver", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 3, name: "count", kind: "scalar", T: 13 /* ScalarType.UINT32 */ },
  ],
);

/**
 * @generated from message io.defang.v1.Resource
 */
export const Resource = /*@__PURE__*/ proto3.makeMessageType(
  "io.defang.v1.Resource",
  () => [
    { no: 1, name: "memory", kind: "scalar", T: 2 /* ScalarType.FLOAT */ },
    { no: 2, name: "cpus", kind: "scalar", T: 2 /* ScalarType.FLOAT */ },
    { no: 3, name: "devices", kind: "message", T: Device, repeated: true },
  ],
);

/**
 * @generated from message io.defang.v1.Resources
 */
export const Resources = /*@__PURE__*/ proto3.makeMessageType(
  "io.defang.v1.Resources",
  () => [
    { no: 1, name: "reservations", kind: "message", T: Resource },
  ],
);

/**
 * @generated from message io.defang.v1.Deploy
 */
export const Deploy = /*@__PURE__*/ proto3.makeMessageType(
  "io.defang.v1.Deploy",
  () => [
    { no: 1, name: "replicas", kind: "scalar", T: 13 /* ScalarType.UINT32 */ },
    { no: 2, name: "resources", kind: "message", T: Resources },
  ],
);

/**
 * @generated from message io.defang.v1.Port
 */
export const Port = /*@__PURE__*/ proto3.makeMessageType(
  "io.defang.v1.Port",
  () => [
    { no: 1, name: "target", kind: "scalar", T: 13 /* ScalarType.UINT32 */ },
    { no: 2, name: "protocol", kind: "enum", T: proto3.getEnumType(Protocol) },
    { no: 3, name: "mode", kind: "enum", T: proto3.getEnumType(Mode) },
  ],
);

/**
 * @generated from message io.defang.v1.Secret
 */
export const Secret = /*@__PURE__*/ proto3.makeMessageType(
  "io.defang.v1.Secret",
  () => [
    { no: 1, name: "source", kind: "scalar", T: 9 /* ScalarType.STRING */ },
  ],
);

/**
 * @generated from message io.defang.v1.Build
 */
export const Build = /*@__PURE__*/ proto3.makeMessageType(
  "io.defang.v1.Build",
  () => [
    { no: 1, name: "context", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 2, name: "dockerfile", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 3, name: "args", kind: "map", K: 9 /* ScalarType.STRING */, V: {kind: "scalar", T: 9 /* ScalarType.STRING */} },
    { no: 4, name: "shm_size", kind: "scalar", T: 2 /* ScalarType.FLOAT */ },
    { no: 5, name: "target", kind: "scalar", T: 9 /* ScalarType.STRING */ },
  ],
);

/**
 * @generated from message io.defang.v1.HealthCheck
 */
export const HealthCheck = /*@__PURE__*/ proto3.makeMessageType(
  "io.defang.v1.HealthCheck",
  () => [
    { no: 1, name: "test", kind: "scalar", T: 9 /* ScalarType.STRING */, repeated: true },
    { no: 2, name: "interval", kind: "scalar", T: 13 /* ScalarType.UINT32 */ },
    { no: 3, name: "timeout", kind: "scalar", T: 13 /* ScalarType.UINT32 */ },
    { no: 4, name: "retries", kind: "scalar", T: 13 /* ScalarType.UINT32 */ },
  ],
);

/**
 * @generated from message io.defang.v1.Service
 */
export const Service = /*@__PURE__*/ proto3.makeMessageType(
  "io.defang.v1.Service",
  () => [
    { no: 1, name: "name", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 2, name: "image", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 3, name: "platform", kind: "enum", T: proto3.getEnumType(Platform) },
    { no: 4, name: "internal", kind: "scalar", T: 8 /* ScalarType.BOOL */ },
    { no: 5, name: "deploy", kind: "message", T: Deploy },
    { no: 6, name: "ports", kind: "message", T: Port, repeated: true },
    { no: 7, name: "environment", kind: "map", K: 9 /* ScalarType.STRING */, V: {kind: "scalar", T: 9 /* ScalarType.STRING */} },
    { no: 8, name: "build", kind: "message", T: Build },
    { no: 9, name: "secrets", kind: "message", T: Secret, repeated: true },
    { no: 10, name: "healthcheck", kind: "message", T: HealthCheck },
    { no: 11, name: "command", kind: "scalar", T: 9 /* ScalarType.STRING */, repeated: true },
    { no: 12, name: "domainname", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 13, name: "init", kind: "scalar", T: 8 /* ScalarType.BOOL */ },
    { no: 14, name: "dns_role", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 15, name: "static_files", kind: "message", T: StaticFiles },
    { no: 16, name: "networks", kind: "enum", T: proto3.getEnumType(Network) },
    { no: 18, name: "redis", kind: "message", T: Redis },
  ],
);

/**
 * @generated from message io.defang.v1.StaticFiles
 */
export const StaticFiles = /*@__PURE__*/ proto3.makeMessageType(
  "io.defang.v1.StaticFiles",
  () => [
    { no: 1, name: "folder", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 2, name: "redirects", kind: "scalar", T: 9 /* ScalarType.STRING */, repeated: true },
  ],
);

/**
 * @generated from message io.defang.v1.Redis
 */
export const Redis = /*@__PURE__*/ proto3.makeMessageType(
  "io.defang.v1.Redis",
  [],
);

/**
 * @generated from message io.defang.v1.Event
 */
export const Event = /*@__PURE__*/ proto3.makeMessageType(
  "io.defang.v1.Event",
  () => [
    { no: 1, name: "specversion", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 2, name: "type", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 3, name: "source", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 4, name: "id", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 5, name: "datacontenttype", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 6, name: "dataschema", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 7, name: "subject", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 8, name: "time", kind: "message", T: Timestamp },
    { no: 9, name: "data", kind: "scalar", T: 12 /* ScalarType.BYTES */ },
  ],
);

/**
 * @generated from message io.defang.v1.PublishRequest
 */
export const PublishRequest = /*@__PURE__*/ proto3.makeMessageType(
  "io.defang.v1.PublishRequest",
  () => [
    { no: 1, name: "event", kind: "message", T: Event },
  ],
);

/**
 * @generated from message io.defang.v1.SubscribeRequest
 */
export const SubscribeRequest = /*@__PURE__*/ proto3.makeMessageType(
  "io.defang.v1.SubscribeRequest",
  () => [
    { no: 1, name: "services", kind: "scalar", T: 9 /* ScalarType.STRING */, repeated: true },
  ],
);

/**
 * @generated from message io.defang.v1.SubscribeResponse
 */
export const SubscribeResponse = /*@__PURE__*/ proto3.makeMessageType(
  "io.defang.v1.SubscribeResponse",
  () => [
    { no: 1, name: "service", kind: "message", T: ServiceInfo },
  ],
);

/**
 * @generated from message io.defang.v1.DelegateSubdomainZoneRequest
 */
export const DelegateSubdomainZoneRequest = /*@__PURE__*/ proto3.makeMessageType(
  "io.defang.v1.DelegateSubdomainZoneRequest",
  () => [
    { no: 1, name: "name_server_records", kind: "scalar", T: 9 /* ScalarType.STRING */, repeated: true },
  ],
);

/**
 * @generated from message io.defang.v1.DelegateSubdomainZoneResponse
 */
export const DelegateSubdomainZoneResponse = /*@__PURE__*/ proto3.makeMessageType(
  "io.defang.v1.DelegateSubdomainZoneResponse",
  () => [
    { no: 1, name: "zone", kind: "scalar", T: 9 /* ScalarType.STRING */ },
  ],
);

/**
 * @generated from message io.defang.v1.WhoAmIResponse
 */
export const WhoAmIResponse = /*@__PURE__*/ proto3.makeMessageType(
  "io.defang.v1.WhoAmIResponse",
  () => [
    { no: 1, name: "tenant", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 2, name: "account", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 3, name: "region", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 4, name: "user_id", kind: "scalar", T: 9 /* ScalarType.STRING */ },
  ],
);

