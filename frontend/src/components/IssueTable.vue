<template>
  <BBTable
    :columnList="columnList"
    :sectionDataSource="issueSectionList"
    :showHeader="true"
    :leftBordered="leftBordered"
    :rightBordered="rightBordered"
    :topBordered="topBordered"
    :bottomBordered="bottomBordered"
    @click-row="clickIssue"
  >
    <template v-slot:header>
      <BBTableHeaderCell
        :leftPadding="4"
        class="w-4 table-cell"
        :title="columnList[0].title"
      />
      <template v-if="mode == 'ALL'">
        <BBTableHeaderCell
          class="w-12 table-cell"
          :title="columnList[1].title"
        />
        <BBTableHeaderCell
          class="w-48 table-cell"
          :title="columnList[2].title"
        />
        <BBTableHeaderCell
          class="w-12 table-cell"
          :title="columnList[3].title"
        />
        <BBTableHeaderCell
          class="w-12 table-cell"
          :title="columnList[4].title"
        />
        <BBTableHeaderCell
          class="w-24 hidden sm:table-cell"
          :title="columnList[5].title"
        />
        <BBTableHeaderCell
          class="w-24 hidden md:table-cell"
          :title="columnList[6].title"
        />
        <BBTableHeaderCell
          class="w-36 hidden sm:table-cell"
          :title="columnList[7].title"
        />
      </template>
      <template v-else-if="mode == 'PROJECT'">
        <BBTableHeaderCell
          class="w-48 table-cell"
          :title="columnList[1].title"
        />
        <BBTableHeaderCell
          class="w-12 table-cell"
          :title="columnList[2].title"
        />
        <BBTableHeaderCell
          class="w-12 table-cell"
          :title="columnList[3].title"
        />
        <BBTableHeaderCell
          class="w-24 hidden sm:table-cell"
          :title="columnList[4].title"
        />
        <BBTableHeaderCell
          class="w-24 hidden md:table-cell"
          :title="columnList[5].title"
        />
        <BBTableHeaderCell
          class="w-36 hidden sm:table-cell"
          :title="columnList[6].title"
        />
      </template>
    </template>
    <template v-slot:body="{ rowData: issue }">
      <BBTableCell :leftPadding="4" class="table-cell">
        <IssueStatusIcon
          :issueStatus="issue.status"
          :taskStatus="activeTask(issue.pipeline).status"
        />
      </BBTableCell>
      <BBTableCell v-if="mode == 'ALL'" class="table-cell text-gray-500">
        <span class="">{{ issue.project.key }}</span>
      </BBTableCell>
      <BBTableCell class="truncate">
        {{ issue.name }}
      </BBTableCell>
      <BBTableCell class="table-cell">
        {{ activeEnvironmentName(issue) }}
      </BBTableCell>
      <BBTableCell class="table-cell">
        {{ activeDatabaseName(issue) }}
      </BBTableCell>
      <BBTableCell class="hidden sm:table-cell">
        <BBStepBar
          :stepList="taskStepList(issue)"
          @click-step="
            (step) => {
              clickIssueStep(issue, step);
            }
          "
        />
      </BBTableCell>
      <BBTableCell class="hidden md:table-cell">
        {{ humanizeTs(issue.updatedTs) }}
      </BBTableCell>
      <BBTableCell class="hidden sm:table-cell">
        <div class="flex flex-row items-center">
          <BBAvatar
            :size="'SMALL'"
            :username="issue.assignee ? issue.assignee.name : 'Unassigned'"
          />
          <span class="ml-2">{{
            issue.assignee ? issue.assignee.name : "Unassigned"
          }}</span>
        </div>
      </BBTableCell>
    </template>
  </BBTable>
</template>

<script lang="ts">
import { reactive, PropType } from "vue";
import { useRouter } from "vue-router";
import {
  BBTableColumn,
  BBTableSectionDataSource,
  BBStep,
  BBStepStatus,
} from "../bbkit/types";
import IssueStatusIcon from "../components/IssueStatusIcon.vue";
import {
  issueSlug,
  activeEnvironment,
  activeDatabase,
  activeTask,
  allTaskList,
  stageSlug,
} from "../utils";
import { Issue, Task } from "../types";

type Mode = "ALL" | "PROJECT";

const columnListMap: Map<Mode, BBTableColumn[]> = new Map([
  [
    "ALL",
    [
      {
        title: "Status",
      },
      {
        title: "Project",
      },
      {
        title: "Name",
      },
      {
        title: "Environment",
      },
      {
        title: "DB",
      },
      {
        title: "Progress",
      },
      {
        title: "Updated",
      },
      {
        title: "Assignee",
      },
    ],
  ],
  [
    "PROJECT",
    [
      {
        title: "Status",
      },
      {
        title: "Name",
      },
      {
        title: "Environment",
      },
      {
        title: "DB",
      },
      {
        title: "Progress",
      },
      {
        title: "Updated",
      },
      {
        title: "Assignee",
      },
    ],
  ],
]);

interface LocalState {
  dataSource: Object[];
}

export default {
  name: "IssueTable",
  components: { IssueStatusIcon },
  props: {
    issueSectionList: {
      required: true,
      type: Object as PropType<BBTableSectionDataSource<Issue>[]>,
    },
    mode: {
      default: "ALL",
      type: String as PropType<Mode>,
    },
    leftBordered: {
      default: true,
      type: Boolean,
    },
    rightBordered: {
      default: true,
      type: Boolean,
    },
    topBordered: {
      default: true,
      type: Boolean,
    },
    bottomBordered: {
      default: true,
      type: Boolean,
    },
  },
  setup(props, ctx) {
    const router = useRouter();

    const state = reactive<LocalState>({
      dataSource: [],
    });

    const activeEnvironmentName = function (issue: Issue) {
      return activeEnvironment(issue.pipeline).name;
    };

    const activeDatabaseName = function (issue: Issue) {
      return activeDatabase(issue.pipeline)?.name;
    };

    const taskStepList = function (issue: Issue): BBStep[] {
      const list: Task[] = allTaskList(issue.pipeline);
      return list.map((task) => {
        let status: BBStepStatus = task.status;
        if (status == "PENDING" || status == "PENDING_APPROVAL") {
          if (activeTask(issue.pipeline).id == task.id) {
            status =
              status == "PENDING"
                ? "PENDING_ACTIVE"
                : "PENDING_APPROVAL_ACTIVE";
          }
        }
        return {
          title: task.name,
          status,
          payload: task,
        };
      });
    };

    const clickIssue = (section: number, row: number) => {
      const issue = props.issueSectionList[section].list[row];
      router.push(`/issue/${issueSlug(issue.name, issue.id)}`);
    };

    const clickIssueStep = (issue: Issue, step: BBStep) => {
      const task = step.payload as Task;
      const stageIndex = issue.pipeline.stageList.findIndex((item, index) => {
        return item.id == task.stage.id;
      });

      router.push({
        name: "workspace.issue.detail",
        params: {
          issueSlug: issueSlug(issue.name, issue.id),
        },
        query: { stage: stageSlug(task.stage.name, stageIndex) },
      });
    };

    return {
      state,
      columnList: columnListMap.get(props.mode),
      activeEnvironmentName,
      activeDatabaseName,
      taskStepList,
      activeTask,
      clickIssue,
      clickIssueStep,
    };
  },
};
</script>
