<template>
  <div class="mx-4 space-y-4 w-160">
    <template v-if="projectId">
      <div v-if="state.project.workflowType == 'VCS'" class="textlabel">
        This project has version control enabled and selecting database below
        will navigate you to the corresponding Git repository to create schema
        change.
      </div>
    </template>
    <template v-else>
      <div class="flex flex-row space-x-2">
        <svg
          class="w-8 h-8 text-control -mt-1.5"
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
          xmlns="http://www.w3.org/2000/svg"
        >
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10"
          ></path>
        </svg>
        <p class="textlabel">
          indicates the project has enabled version control. Selecting database
          belonging to such project will navigate you to the corresponding Git
          repository to create schema change.
        </p>
      </div>
    </template>

    <div
      v-if="projectId && state.project.workflowType == 'UI'"
      class="mt-2 textlabel"
    >
      <div class="radio-set-row">
        <div class="radio">
          <input
            tabindex="-1"
            type="radio"
            class="btn"
            value="SINGLE_DB"
            v-model="state.alterType"
          />
          <label class="label"> Alter single database </label>
        </div>
        <div class="radio">
          <input
            tabindex="-1"
            type="radio"
            class="btn"
            value="MULTI_DB"
            v-model="state.alterType"
          />
          <label class="label"> Alter multiple databases </label>
        </div>
      </div>
    </div>

    <template v-if="projectId && state.alterType == 'MULTI_DB'">
      <div class="textinfolabel">
        For each environment, your can select a database to alter its schema or
        just skip that environment. This allows you to compose a single pipeline
        to propagate schema changes across multiple environments.
      </div>
      <div class="space-y-4">
        <div v-for="(environment, index) in environmentList" :key="index">
          <div class="mb-2">{{ environment.name }}</div>
          <div class="relative bg-white rounded-md -space-y-px">
            <template
              v-for="(database, index) in databaseList.filter(
                (item) => item.instance.environment.id == environment.id
              )"
              :key="index"
            >
              <label
                class="
                  border-control-border
                  relative
                  border
                  p-3
                  flex flex-col
                  md:pl-4 md:pr-6 md:grid md:grid-cols-2
                "
                :class="
                  database.syncStatus == 'OK'
                    ? 'cursor-pointer'
                    : 'cursor-not-allowed'
                "
              >
                <div class="radio text-sm">
                  <input
                    v-if="database.syncStatus == 'OK'"
                    type="radio"
                    class="btn"
                    :checked="
                      state.selectedDatabaseIdForEnvironment.get(
                        environment.id
                      ) == database.id
                    "
                    @change="
                      selectDatabaseIdForEnvironment(
                        database.id,
                        environment.id
                      )
                    "
                  />
                  <span
                    class="font-medium"
                    :class="
                      database.syncStatus == 'OK'
                        ? 'ml-2 text-main'
                        : 'ml-6 text-control-light'
                    "
                    >{{ database.name }}</span
                  >
                </div>
                <p
                  class="
                    textinfolabel
                    ml-6
                    pl-1
                    text-sm
                    md:ml-0 md:pl-0 md:text-right
                  "
                >
                  Last sync status:
                  <span
                    :class="
                      database.syncStatus == 'OK'
                        ? 'textlabel'
                        : 'text-sm font-medium text-error'
                    "
                    >{{ database.syncStatus }}</span
                  >
                </p>
              </label>
            </template>
            <label
              class="
                border-control-border
                relative
                border
                p-3
                flex flex-col
                cursor-pointer
                md:pl-4 md:pr-6 md:grid md:grid-cols-3
              "
            >
              <div class="radio space-x-2 text-sm">
                <input
                  type="radio"
                  class="btn"
                  :checked="
                    state.selectedDatabaseIdForEnvironment.get(environment.id)
                      ? 0
                      : 1
                  "
                  @input="clearDatabaseIdForEnvironment(environment.id)"
                />
                <span class="ml-3 font-medium text-main">SKIP</span>
              </div>
            </label>
          </div>
        </div>
      </div>
    </template>
    <template v-else>
      <DatabaseTable
        :mode="projectId ? 'PROJECT_SHORT' : 'ALL_SHORT'"
        :bordered="true"
        :customClick="true"
        :databaseList="databaseList"
        @select-database-id="selectDatabaseId"
      />
    </template>
    <!-- Create button group -->
    <div class="pt-4 border-t border-block-border flex justify-end">
      <button
        type="button"
        class="btn-normal py-2 px-4"
        @click.prevent="cancel"
      >
        Cancel
      </button>
      <button
        v-if="state.alterType == 'MULTI_DB'"
        class="btn-primary ml-3 inline-flex justify-center py-2 px-4"
        :disabled="!allowGenerateMultiDb"
        @click.prevent="generateMultDb"
      >
        Next
      </button>
    </div>
  </div>
</template>

<script lang="ts">
import { computed, reactive, onMounted, onUnmounted, PropType } from "vue";
import { useStore } from "vuex";
import { useRouter } from "vue-router";
import DatabaseTable from "../components/DatabaseTable.vue";
import {
  baseDirectoryWebURL,
  DatabaseId,
  EnvironmentId,
  Project,
  ProjectId,
  Repository,
} from "../types";
import { sortDatabaseList } from "../utils";
import { cloneDeep } from "lodash";

type AlterType = "SINGLE_DB" | "MULTI_DB";

interface LocalState {
  project?: Project;
  alterType: AlterType;
  selectedDatabaseIdForEnvironment: Map<EnvironmentId, DatabaseId>;
}

export default {
  name: "AlterSchemaPrepForm",
  emits: ["dismiss"],
  props: {
    projectId: {
      type: Number as PropType<ProjectId>,
    },
  },
  components: {
    DatabaseTable,
  },
  setup(props, { emit }) {
    const store = useStore();
    const router = useRouter();

    const currentUser = computed(() => store.getters["auth/currentUser"]());

    const keyboardHandler = (e: KeyboardEvent) => {
      if (e.code == "Escape") {
        cancel();
      }
    };

    onMounted(() => {
      document.addEventListener("keydown", keyboardHandler);
    });

    onUnmounted(() => {
      document.removeEventListener("keydown", keyboardHandler);
    });

    const state = reactive<LocalState>({
      project: props.projectId
        ? store.getters["project/projectById"](props.projectId)
        : undefined,
      alterType: "SINGLE_DB",
      selectedDatabaseIdForEnvironment: new Map(),
    });

    const environmentList = computed(() => {
      return store.getters["environment/environmentList"](["NORMAL"]);
    });

    const databaseList = computed(() => {
      var list;
      if (props.projectId) {
        list = store.getters["database/databaseListByProjectId"](
          props.projectId
        );
      } else {
        list = store.getters["database/databaseListByPrincipalId"](
          currentUser.value.id
        );
      }

      return sortDatabaseList(cloneDeep(list), environmentList.value);
    });

    const allowGenerateMultiDb = computed(() => {
      return state.selectedDatabaseIdForEnvironment.size > 0;
    });

    const generateMultDb = () => {
      const databaseIdList: DatabaseId[] = [];
      for (var i = 0; i < environmentList.value.length; i++) {
        if (
          state.selectedDatabaseIdForEnvironment.get(
            environmentList.value[i].id
          )
        ) {
          databaseIdList.push(
            state.selectedDatabaseIdForEnvironment.get(
              environmentList.value[i].id
            )!
          );
        }
      }
      router.push({
        name: "workspace.issue.detail",
        params: {
          issueSlug: "new",
        },
        query: {
          template: "bb.issue.database.schema.update",
          name: `Alter schema`,
          project: props.projectId,
          databaseList: databaseIdList.join(","),
        },
      });
    };

    const selectDatabaseId = (databaseId: DatabaseId) => {
      emit("dismiss");

      const database = store.getters["database/databaseById"](databaseId);

      if (database.project.workflowType == "UI") {
        router.push({
          name: "workspace.issue.detail",
          params: {
            issueSlug: "new",
          },
          query: {
            template: "bb.issue.database.schema.update",
            name: `[${database.name}] Alter schema`,
            project: database.project.id,
            databaseList: database.id,
          },
        });
      } else if (database.project.workflowType == "VCS") {
        store
          .dispatch(
            "repository/fetchRepositoryByProjectId",
            database.project.id
          )
          .then((repository: Repository) => {
            window.open(baseDirectoryWebURL(repository), "_blank");
          });
      }
    };

    const selectDatabaseIdForEnvironment = (
      databaseId: DatabaseId,
      environmentId: EnvironmentId
    ) => {
      state.selectedDatabaseIdForEnvironment.set(environmentId, databaseId);
    };

    const clearDatabaseIdForEnvironment = (environmentId: EnvironmentId) => {
      state.selectedDatabaseIdForEnvironment.delete(environmentId);
    };

    const cancel = () => {
      emit("dismiss");
    };

    return {
      state,
      environmentList,
      databaseList,
      allowGenerateMultiDb,
      generateMultDb,
      selectDatabaseId,
      selectDatabaseIdForEnvironment,
      clearDatabaseIdForEnvironment,
      cancel,
    };
  },
};
</script>
