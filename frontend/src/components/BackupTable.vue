<template>
  <BBTable
    :columnList="columnList"
    :sectionDataSource="backupSectionList"
    :showHeader="true"
    :rowClickable="false"
    :leftBordered="true"
    :rightBordered="true"
  >
    <template v-slot:header>
      <BBTableHeaderCell
        :leftPadding="4"
        class="w-4"
        :title="columnList[0].title"
      />
      <BBTableHeaderCell class="w-16" :title="columnList[1].title" />
      <BBTableHeaderCell class="w-48" :title="columnList[2].title" />
      <BBTableHeaderCell class="w-16" :title="columnList[3].title" />
      <BBTableHeaderCell class="w-16" :title="columnList[4].title" />
      <BBTableHeaderCell
        v-if="allowEdit"
        class="w-4"
        :title="columnList[5].title"
      />
    </template>
    <template v-slot:body="{ rowData: backup }">
      <BBTableCell :leftPadding="4">
        <span
          class="flex items-center justify-center rounded-full select-none"
          :class="statusIconClass(backup)"
        >
          <template v-if="backup.status == 'PENDING_CREATE'">
            <span
              class="h-2 w-2 bg-blue-600 hover:bg-blue-700 rounded-full"
              style="
                animation: pulse 2.5s cubic-bezier(0.4, 0, 0.6, 1) infinite;
              "
            >
            </span>
          </template>
          <template v-else-if="backup.status == 'DONE'">
            <svg
              class="w-4 h-4"
              xmlns="http://www.w3.org/2000/svg"
              viewBox="0 0 20 20"
              fill="currentColor"
              aria-hidden="true"
            >
              <path
                fill-rule="evenodd"
                d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z"
                clip-rule="evenodd"
              />
            </svg>
          </template>
          <template v-else-if="backup.status == 'FAILED'">
            <span
              class="
                h-2
                w-2
                rounded-full
                text-center
                pb-6
                font-normal
                text-base
              "
              aria-hidden="true"
              >!</span
            >
          </template>
        </span>
      </BBTableCell>
      <BBTableCell>
        {{ backup.name }}
      </BBTableCell>
      <BBTableCell class="tooltip-wrapper">
        <span v-if="backup.comment.length > 100" class="tooltip">{{
          backup.comment
        }}</span>
        {{
          backup.comment.length > 100
            ? backup.comment.substring(0, 100) + "..."
            : backup.comment
        }}
      </BBTableCell>
      <BBTableCell>
        {{ humanizeTs(backup.createdTs) }}
      </BBTableCell>
      <BBTableCell>
        {{ backup.creator.name }}
      </BBTableCell>
      <BBTableCell v-if="allowEdit">
        <button
          class="normal-link"
          @click.stop="
            () => {
              state.restoredBackup = backup;
              state.showRestoreBackupModal = true;
            }
          "
        >
          Restore
        </button>
      </BBTableCell>
    </template>
  </BBTable>
  <BBModal
    v-if="state.showRestoreBackupModal"
    :title="`Restore backup '${state.restoredBackup.name}' to a new database`"
    @close="
      () => {
        state.showRestoreBackupModal = false;
        state.restoredBackup = undefined;
      }
    "
  >
    <CreateDatabasePrepForm
      :projectId="database.project.id"
      :environmentId="database.instance.environment.id"
      :instanceId="database.instance.id"
      :backup="state.restoredBackup"
      @dismiss="
        () => {
          state.showRestoreBackupModal = false;
          state.restoredBackup = undefined;
        }
      "
    />
  </BBModal>
</template>

<script lang="ts">
import { computed, PropType, reactive } from "vue";
import { BBTableColumn, BBTableSectionDataSource } from "../bbkit/types";
import { Backup, Database } from "../types";
import { bytesToString } from "../utils";
import { useStore } from "vuex";
import CreateDatabasePrepForm from "../components/CreateDatabasePrepForm.vue";

const EDIT_COLUMN_LIST: BBTableColumn[] = [
  {
    title: "Status",
  },
  {
    title: "Name",
  },
  {
    title: "Comment",
  },
  {
    title: "Time",
  },
  {
    title: "Creator",
  },
  {
    title: "",
  },
];

const NON_EDIT_COLUMN_LIST: BBTableColumn[] = [
  {
    title: "Status",
  },
  {
    title: "Name",
  },
  {
    title: "Comment",
  },
  {
    title: "Time",
  },
  {
    title: "Creator",
  },
];

interface LocalState {
  showRestoreBackupModal: boolean;
  restoredBackup?: Backup;
}

export default {
  name: "BackupTable",
  components: { CreateDatabasePrepForm },
  props: {
    database: {
      required: true,
      type: Object as PropType<Database>,
    },
    backupList: {
      required: true,
      type: Object as PropType<Backup[]>,
    },
    allowEdit: {
      required: true,
      type: Boolean,
    },
  },
  setup(props, ctx) {
    const store = useStore();

    const state = reactive<LocalState>({
      showRestoreBackupModal: false,
    });

    const backupSectionList = computed(() => {
      const manualList: Backup[] = [];
      const automaticList: Backup[] = [];
      const sectionList: BBTableSectionDataSource<Backup>[] = [
        {
          title: "Manual",
          list: manualList,
        },
        {
          title: "Automatic",
          list: automaticList,
        },
      ];

      for (const backup of props.backupList) {
        if (backup.type == "MANUAL") {
          manualList.push(backup);
        } else if (backup.type == "AUTOMATIC") {
          automaticList.push(backup);
        }
      }

      return sectionList;
    });

    const statusIconClass = (backup: Backup) => {
      let iconClass = "w-5 h-5";
      switch (backup.status) {
        case "PENDING_CREATE":
          return (
            iconClass +
            " bg-white border-2 border-blue-600 text-blue-600 hover:text-blue-700 hover:border-blue-700"
          );
        case "DONE":
          return iconClass + " bg-success hover:bg-success-hover text-white";
        case "FAILED":
          return (
            iconClass +
            " bg-error text-white hover:text-white hover:bg-error-hover"
          );
      }
    };

    const columnList = computed(() => {
      return props.allowEdit ? EDIT_COLUMN_LIST : NON_EDIT_COLUMN_LIST;
    });

    return {
      state,
      columnList,
      bytesToString,
      backupSectionList,
      statusIconClass,
    };
  },
};
</script>
