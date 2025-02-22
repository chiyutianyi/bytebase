import { RowStatus } from "./common";
import { MemberId, PrincipalId, ProjectId } from "./id";
import { OAuthToken } from "./oauth";
import { Principal } from "./principal";
import { ExternalRepositoryInfo, RepositoryConfig } from "./repository";
import { VCS } from "./vcs";

export type ProjectRoleType = "OWNER" | "DEVELOPER";

export type ProjectWorkflowType = "UI" | "VCS";

export type ProjectVisibility = "PUBLIC" | "PRIVATE";

// Project
export type Project = {
  id: ProjectId;

  // Standard fields
  creator: Principal;
  updater: Principal;
  createdTs: number;
  updatedTs: number;
  rowStatus: RowStatus;

  // Domain specific fields
  name: string;
  key: string;
  // Returns the member list directly because we need it quite frequently in order
  // to do various access check.
  memberList: ProjectMember[];
  workflowType: ProjectWorkflowType;
  visibility: ProjectVisibility;
};

export type ProjectCreate = {
  // Domain specific fields
  name: string;
  key: string;
};

export type ProjectPatch = {
  // Standard fields
  rowStatus?: RowStatus;

  // Domain specific fields
  name?: string;
  key?: string;
};

// Project Member
export type ProjectMember = {
  id: MemberId;

  // Related fields
  project: Project;

  // Standard fields
  creator: Principal;
  createdTs: number;
  updater: Principal;
  updatedTs: number;

  // Domain specific fields
  role: ProjectRoleType;
  principal: Principal;
};

export type ProjectMemberCreate = {
  // Domain specific fields
  principalId: PrincipalId;
  role: ProjectRoleType;
};

export type ProjectMemberPatch = {
  // Domain specific fields
  role: ProjectRoleType;
};

export type ProjectRepositoryConfig = {
  vcs: VCS;
  code: string;
  token: OAuthToken;
  repositoryInfo: ExternalRepositoryInfo;
  repositoryConfig: RepositoryConfig;
};
