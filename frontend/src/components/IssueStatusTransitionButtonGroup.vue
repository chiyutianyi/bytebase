<template>
  <template v-if="create">
    <button
      type="button"
      class="btn-primary px-4 py-2"
      @click.prevent="doCreate"
      :disabled="!allowCreate"
    >
      Create
    </button>
  </template>
  <template v-else>
    <div
      v-if="applicableTaskStatusTransitionList.length > 0"
      class="flex space-x-2"
    >
      <template
        v-for="(transition, index) in applicableTaskStatusTransitionList"
        :key="index"
      >
        <button
          type="button"
          :class="transition.buttonClass"
          @click.prevent="tryStartTaskStatusTransition(transition)"
        >
          {{ transition.buttonName }}
        </button>
      </template>
      <template v-if="applicableIssueStatusTransitionList.length > 0">
        <button
          type="button"
          @click.prevent="$refs.menu.toggle($event)"
          @contextmenu.capture.prevent="$refs.menu.toggle($event)"
          class="text-control-light"
          id="user-menu"
          aria-label="User menu"
          aria-haspopup="true"
        >
          <svg
            class="w-6 h-6"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
            xmlns="http://www.w3.org/2000/svg"
          >
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d="M12 5v.01M12 12v.01M12 19v.01M12 6a1 1 0 110-2 1 1 0 010 2zm0 7a1 1 0 110-2 1 1 0 010 2zm0 7a1 1 0 110-2 1 1 0 010 2z"
            ></path>
          </svg>
        </button>
        <BBContextMenu ref="menu" class="origin-top-right mt-10 w-42">
          <template
            v-for="(transition, index) in applicableIssueStatusTransitionList"
            :key="index"
          >
            <div
              @click.prevent="tryStartIssueStatusTransition(transition)"
              class="menu-item"
              role="menuitem"
            >
              {{ transition.buttonName }}
            </div>
          </template>
        </BBContextMenu>
      </template>
    </div>
    <template v-else>
      <div
        if="applicableIssueStatusTransitionList.length > 0"
        class="flex space-x-2"
      >
        <template
          v-for="(transition, index) in applicableIssueStatusTransitionList"
          :key="index"
        >
          <button
            type="button"
            :class="transition.buttonClass"
            :disabled="!allowIssueStatusTransition(transition)"
            @click.prevent="tryStartIssueStatusTransition(transition)"
          >
            {{ transition.buttonName }}
          </button>
        </template>
      </div>
    </template>
  </template>
  <BBModal
    v-if="updateStatusModalState.show"
    :title="updateStatusModalState.title"
    @close="updateStatusModalState.show = false"
  >
    <StatusTransitionForm
      :mode="updateStatusModalState.mode"
      :okText="updateStatusModalState.okText"
      :issue="issue"
      :transition="updateStatusModalState.transition"
      :outputFieldList="issueTemplate.outputFieldList"
      @submit="
        (comment) => {
          updateStatusModalState.show = false;
          if (updateStatusModalState.mode == 'ISSUE') {
            doIssueStatusTransition(updateStatusModalState.transition, comment);
          } else if (updateStatusModalState.mode == 'TASK') {
            doTaskStatusTransition(
              updateStatusModalState.transition,
              updateStatusModalState.payload,
              comment
            );
          }
        }
      "
      @cancel="
        () => {
          updateStatusModalState.show = false;
        }
      "
    />
  </BBModal>
</template>

<script lang="ts">
import { PropType, computed, reactive } from "vue";
import { Store, useStore } from "vuex";
import StatusTransitionForm from "../components/StatusTransitionForm.vue";
import {
  activeTask,
  allTaskList,
  applicableTaskTransition,
  isDBAOrOwner,
  TaskStatusTransition,
} from "../utils";
import {
  ASSIGNEE_APPLICABLE_ACTION_LIST,
  CREATOR_APPLICABLE_ACTION_LIST,
  Issue,
  IssueCreate,
  IssueStatusTransition,
  IssueStatusTransitionType,
  ISSUE_STATUS_TRANSITION_LIST,
  ONBOARDING_ISSUE_ID,
  Principal,
  SYSTEM_BOT_ID,
  Task,
  TaskId,
} from "../types";
import { IssueTemplate } from "../plugins";
import isEmpty from "lodash-es/isEmpty";

interface UpdateStatusModalState {
  mode: "ISSUE" | "TASK";
  show: boolean;
  style: string;
  okText: string;
  title: string;
  transition?: IssueStatusTransition | TaskStatusTransition;
  payload?: Task;
}

export type IssueContext = {
  store: Store<any>;
  currentUser: Principal;
  create: boolean;
  issue: Issue | IssueCreate;
};

export default {
  name: "IssueStatusTransitionButtonGroup",
  emits: ["create", "change-issue-status", "change-task-status"],
  props: {
    create: {
      required: true,
      type: Boolean,
    },
    issue: {
      required: true,
      type: Object as PropType<Issue | IssueCreate>,
    },
    issueTemplate: {
      required: true,
      type: Object as PropType<IssueTemplate>,
    },
  },
  components: {
    StatusTransitionForm,
  },
  setup(props, { emit }) {
    const store = useStore();

    const updateStatusModalState = reactive<UpdateStatusModalState>({
      mode: "ISSUE",
      show: false,
      style: "INFO",
      okText: "OK",
      title: "",
    });

    const currentUser = computed(() => store.getters["auth/currentUser"]());

    const issueContext = computed((): IssueContext => {
      return {
        store,
        currentUser: currentUser.value,
        create: props.create,
        issue: props.issue,
      };
    });

    const applicableTaskStatusTransitionList = computed(
      (): TaskStatusTransition[] => {
        if ((props.issue as Issue).id == ONBOARDING_ISSUE_ID) {
          return [];
        }
        switch ((props.issue as Issue).status) {
          case "DONE":
          case "CANCELED":
            return [];
          case "OPEN": {
            let list: TaskStatusTransition[] = [];

            // Allow assignee, or assignee is the system bot and current user is DBA or owner
            if (
              currentUser.value.id === (props.issue as Issue).assignee?.id ||
              ((props.issue as Issue).assignee?.id == SYSTEM_BOT_ID &&
                isDBAOrOwner(currentUser.value.role))
            ) {
              list = applicableTaskTransition((props.issue as Issue).pipeline);
            }

            return list;
          }
        }
      }
    );

    const tryStartTaskStatusTransition = (transition: TaskStatusTransition) => {
      updateStatusModalState.mode = "TASK";
      updateStatusModalState.okText = transition.buttonName;
      const task = activeTask((props.issue as Issue).pipeline);
      switch (transition.type) {
        case "RUN":
          updateStatusModalState.style = "INFO";
          updateStatusModalState.title = `Run '${task.name}'?`;
          break;
        case "APPROVE":
          updateStatusModalState.style = "INFO";
          updateStatusModalState.title = `Approve '${task.name}'?`;
          break;
        case "RETRY":
          updateStatusModalState.style = "INFO";
          updateStatusModalState.title = `Retry '${task.name}'?`;
          break;
        case "CANCEL":
          updateStatusModalState.style = "INFO";
          updateStatusModalState.title = `Cancel '${task.name}'?`;
          break;
        case "SKIP":
          updateStatusModalState.style = "INFO";
          updateStatusModalState.title = `Skip '${task.name}'?`;
          break;
      }
      updateStatusModalState.transition = transition;
      updateStatusModalState.payload = task;
      updateStatusModalState.show = true;
    };

    const doTaskStatusTransition = (
      transition: TaskStatusTransition,
      task: Task,
      comment: string
    ) => {
      emit("change-task-status", task, transition.to, comment);
    };

    const applicableIssueStatusTransitionList = computed(
      (): IssueStatusTransition[] => {
        if ((props.issue as Issue).id == ONBOARDING_ISSUE_ID) {
          return [];
        }
        const list: IssueStatusTransitionType[] = [];
        // Allow assignee, or assignee is the system bot and current user is DBA or owner
        if (
          currentUser.value.id === (props.issue as Issue).assignee?.id ||
          ((props.issue as Issue).assignee?.id == SYSTEM_BOT_ID &&
            isDBAOrOwner(currentUser.value.role))
        ) {
          list.push(
            ...ASSIGNEE_APPLICABLE_ACTION_LIST.get(
              (props.issue as Issue).status
            )!
          );
        }
        if (currentUser.value.id === (props.issue as Issue).creator.id) {
          CREATOR_APPLICABLE_ACTION_LIST.get(
            (props.issue as Issue).status
          )!.forEach((item) => {
            if (list.indexOf(item) == -1) {
              list.push(item);
            }
          });
        }

        return list
          .filter((item) => {
            const pipeline = (props.issue as Issue).pipeline;
            const currentTask = activeTask(pipeline);
            // Disallow any issue status transition if the active task is in RUNNING state.
            if (currentTask.status == "RUNNING") {
              return false;
            }

            const taskList = allTaskList(pipeline);
            // Don't display the Resolve action if the last task is NOT in DONE status.
            if (
              item == "RESOLVE" &&
              taskList.length > 0 &&
              (currentTask.id != taskList[taskList.length - 1].id ||
                currentTask.status != "DONE")
            ) {
              return false;
            }

            return true;
          })
          .map(
            (type: IssueStatusTransitionType) =>
              ISSUE_STATUS_TRANSITION_LIST.get(type)!
          )
          .reverse();
      }
    );

    const allowIssueStatusTransition = (
      transition: IssueStatusTransition
    ): boolean => {
      const issue: Issue = props.issue as Issue;
      if (transition.type == "RESOLVE") {
        // Returns false if any of the required output fields is not provided.
        for (let i = 0; i < props.issueTemplate.outputFieldList.length; i++) {
          const field = props.issueTemplate.outputFieldList[i];
          if (!field.resolved(issueContext.value)) {
            return false;
          }
        }
        return true;
      }
      return true;
    };

    const tryStartIssueStatusTransition = (
      transition: IssueStatusTransition
    ) => {
      updateStatusModalState.mode = "ISSUE";
      updateStatusModalState.okText = transition.buttonName;
      switch (transition.type) {
        case "RESOLVE":
          updateStatusModalState.style = "SUCCESS";
          updateStatusModalState.title = "Resolve issue?";
          break;
        case "CANCEL":
          updateStatusModalState.style = "INFO";
          updateStatusModalState.title = "Cancel this entire issue?";
          break;
        case "REOPEN":
          updateStatusModalState.style = "INFO";
          updateStatusModalState.title = "Reopen issue?";
          break;
      }
      updateStatusModalState.transition = transition;
      updateStatusModalState.show = true;
    };

    const doIssueStatusTransition = (
      transition: IssueStatusTransition,
      comment: string
    ) => {
      emit("change-issue-status", transition.to, comment);
    };

    const allowCreate = computed(() => {
      const newIssue = props.issue as IssueCreate;
      if (isEmpty(newIssue.name)) {
        return false;
      }

      if (!newIssue.assigneeId) {
        return false;
      }

      if (
        newIssue.type == "bb.issue.database.create" ||
        newIssue.type == "bb.issue.database.schema.update"
      ) {
        for (const stage of newIssue.pipeline.stageList) {
          for (const task of stage.taskList) {
            if (
              task.type == "bb.task.database.create" ||
              task.type == "bb.task.database.schema.update"
            ) {
              if (isEmpty(task.statement)) {
                return false;
              }
            }
          }
        }
      }

      for (const field of props.issueTemplate.inputFieldList) {
        if (
          field.type != "Boolean" && // Switch is boolean value which always is present
          !field.resolved(issueContext.value)
        ) {
          return false;
        }
      }
      return true;
    });

    const doCreate = () => {
      emit("create");
    };

    return {
      updateStatusModalState,
      applicableTaskStatusTransitionList,
      tryStartTaskStatusTransition,
      doTaskStatusTransition,
      applicableIssueStatusTransitionList,
      allowIssueStatusTransition,
      tryStartIssueStatusTransition,
      doIssueStatusTransition,
      allowCreate,
      doCreate,
    };
  },
};
</script>
