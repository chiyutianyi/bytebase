import { IssueCreate, UNKNOWN_ID } from "../../types";
import { IssueTemplate, TemplateContext } from "../types";

const template: IssueTemplate = {
  type: "bb.issue.general",
  buildIssue: (
    ctx: TemplateContext
  ): Omit<IssueCreate, "projectId" | "creatorId"> => {
    return {
      name: "",
      type: "bb.issue.general",
      description: "",
      assigneeId: UNKNOWN_ID,
      pipeline: {
        stageList: [
          {
            name: "Troubleshoot database",
            environmentId: ctx.environmentList[0].id,
            taskList: [
              {
                name: "Troubleshoot database",
                status: "PENDING_APPROVAL",
                type: "bb.task.general",
                instanceId: UNKNOWN_ID,
                databaseId: UNKNOWN_ID,
                statement: "",
                rollbackStatement: "",
              },
            ],
          },
        ],
        name: "Create database pipeline",
      },
      payload: {},
    };
  },
  inputFieldList: [],
  outputFieldList: [],
};

export default template;
