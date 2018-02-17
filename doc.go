// Package protokit provides core types used for parsing proto3 documents.
//
// Given the following proto file (api.proto):
//
//     syntax = "proto3";
//     option go_package = "services";
//
//     import "google/protobuf/timestamp.proto";
//
//     // A service for managing "todo" items.
//     //
//     // Add, complete, and remove your items on your todo lists.
//     service Todo {
//       // Creates a new todo list
//       rpc CreateList(CreateListRequest) returns (CreateListResponse);
//     }
//
//     // An enumeration of list types
//     enum ListType {
//       REMINDERS = 0; // The reminders type.
//       CHECKLIST = 1; // The checklist type.
//     }
//
//     // The list object.
//     message List {
//       int64 id                             = 1; // The id of the list.
//       string name                          = 2; // The name of the list.
//       ListType type                        = 3; // The list type.
//       google.protobuf.Timestamp created_at = 4; // The timestamp of the creation.
//     }
//
//     // A request object for creating lists
//     message CreateListRequest {
//       string name   = 1; // The name of the list.
//       ListType type = 2; // The type of list to create.
//     }
//
//     // A response for created lists
//     message CreateListResponse {
//       List list = 1; // The newly created list.
//     }
//
//
// A starter kit for building protoc-plugins. Rather than write your own, you can just use an existing one.
package protokit
