<template>
  <div>
    <p class="text-lg font-medium leading-7 text-main">Manage members</p>
    <div v-if="allowAddMember" class="mt-4 w-full flex justify-start">
      <!-- To prevent jiggling when showing the error text -->
      <div :class="state.error ? 'space-y-1' : 'space-y-6'">
        <div class="space-y-2">
          <div
            class="flex flex-row justify-between py-0.5 select-none space-x-4"
          >
            <div class="w-64">
              <MemberSelect
                id="user"
                name="user"
                :required="false"
                :placeholder="'Select user'"
                :selectedId="state.principalId"
                @select-principal-id="
                  (principalId) => {
                    state.principalId = principalId;
                    clearValidationError();
                    validateMember();
                  }
                "
              />
            </div>
            <div v-if="hasAdminFeature" class="radio-set-row">
              <div class="radio">
                <input
                  :name="`member_role`"
                  tabindex="-1"
                  type="radio"
                  class="btn"
                  value="OWNER"
                  v-model="state.role"
                />
                <label class="label"> Owner </label>
              </div>
              <div class="radio">
                <input
                  :name="`member_role`"
                  tabindex="-1"
                  type="radio"
                  class="btn"
                  value="DEVELOPER"
                  v-model="state.role"
                />
                <label class="label"> Developer </label>
              </div>
            </div>
            <button
              type="button"
              class="btn-primary items-center"
              :disabled="!hasValidMember"
              @click.prevent="addMember"
            >
              <svg
                class="mr-2 w-5 h-5"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
                xmlns="http://www.w3.org/2000/svg"
              >
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="2"
                  d="M18 9v3m0 0v3m0-3h3m-3 0h-3m-2-5a4 4 0 11-8 0 4 4 0 018 0zM3 20a6 6 0 0112 0v1H3v-1z"
                ></path>
              </svg>
              Add member
            </button>
          </div>
        </div>

        <div class="flex justify-start" id="state-error">
          <span v-if="state.error" class="text-sm text-error">
            {{ state.error }}
          </span>
        </div>
      </div>
    </div>
    <ProjectMemberTable :project="project" />
  </div>
</template>

<script lang="ts">
import { computed, PropType, reactive } from "vue";
import { useStore } from "vuex";
import MemberSelect from "../components/MemberSelect.vue";
import ProjectMemberTable from "../components/ProjectMemberTable.vue";
import {
  DEFAULT_PROJECT_ID,
  PrincipalId,
  Project,
  ProjectMember,
  ProjectMemberCreate,
  ProjectRoleType,
  UNKNOWN_ID,
} from "../types";
import { isOwner, isProjectOwner } from "../utils";

interface LocalState {
  principalId: PrincipalId;
  role: ProjectRoleType;
  error: string;
}

export default {
  name: "ProjectMemberPanel",
  components: { MemberSelect, ProjectMemberTable },
  props: {
    project: {
      required: true,
      type: Object as PropType<Project>,
    },
  },
  setup(props, ctx) {
    const store = useStore();

    const currentUser = computed(() => store.getters["auth/currentUser"]());

    const state = reactive<LocalState>({
      principalId: UNKNOWN_ID,
      role: "DEVELOPER",
      error: "",
    });

    const hasAdminFeature = computed(() =>
      store.getters["plan/feature"]("bb.admin")
    );

    const allowAddMember = computed(() => {
      if (props.project.id == DEFAULT_PROJECT_ID) {
        return false;
      }

      if (props.project.rowStatus == "ARCHIVED") {
        return false;
      }

      // Allow workspace owner here in case project owners are not available.
      if (isOwner(currentUser.value.role)) {
        return true;
      }

      for (const member of props.project.memberList) {
        if (member.principal.id == currentUser.value.id) {
          if (isProjectOwner(member.role)) {
            return true;
          }
        }
      }
      return false;
    });

    const hasValidMember = computed(() => {
      return (
        state.principalId != UNKNOWN_ID && validateInviteInternal().length == 0
      );
    });

    const validateInviteInternal = (): string => {
      if (state.principalId != UNKNOWN_ID) {
        if (
          props.project.memberList.find((item: ProjectMember) => {
            return item.principal.id == state.principalId;
          })
        ) {
          return "Already a project member";
        }
      }
      return "";
    };

    const validateMember = () => {
      state.error = validateInviteInternal();
    };

    const clearValidationError = () => {
      state.error = "";
    };

    const addMember = () => {
      // If admin feature is NOT enabled, then we set every member to OWNER role.
      const projectMember: ProjectMemberCreate = {
        principalId: state.principalId,
        role: hasAdminFeature.value ? state.role : "OWNER",
      };
      const member = store.getters["member/memberByPrincipalId"](
        state.principalId
      );
      store
        .dispatch("project/createdMember", {
          projectId: props.project.id,
          projectMember,
        })
        .then(() => {
          store.dispatch("notification/pushNotification", {
            module: "bytebase",
            style: "SUCCESS",
            title: `Successfully added ${member.principal.name} to the project.`,
          });
        });

      state.principalId = UNKNOWN_ID;
      state.role = "DEVELOPER";
      state.error = "";
    };

    return {
      state,
      hasAdminFeature,
      allowAddMember,
      validateMember,
      clearValidationError,
      hasValidMember,
      addMember,
    };
  },
};
</script>
