<template>
  <BBTable
    :columnList="columnList"
    :dataSource="repositoryList"
    :showHeader="true"
    :leftBordered="true"
    :rightBordered="true"
    @click-row="clickRepository"
  >
    <template v-slot:body="{ rowData: repository }">
      <BBTableCell :leftPadding="4" class="w-16">
        {{ projectName(repository.project) }}
      </BBTableCell>
      <BBTableCell class="w-32">
        {{ repository.fullPath }}
      </BBTableCell>
      <BBTableCell class="w-16">
        <div class="flex flex-row items-center">
          <PrincipalAvatar :principal="repository.creator" :size="'SMALL'" />
          <span class="ml-2">{{ repository.creator.name }}</span>
        </div>
      </BBTableCell>
      <BBTableCell class="w-16">
        {{ humanizeTs(repository.createdTs) }}
      </BBTableCell>
    </template>
  </BBTable>
</template>

<script lang="ts">
import { PropType } from "vue";
import { useRouter } from "vue-router";
import PrincipalAvatar from "./PrincipalAvatar.vue";
import { BBTableColumn } from "../bbkit/types";
import { projectSlug } from "../utils";
import { Repository } from "../types";

const columnList: BBTableColumn[] = [
  {
    title: "Project",
  },
  {
    title: "Repository",
  },
  {
    title: "Creator",
  },
  {
    title: "Created",
  },
];

export default {
  name: "RepositoryTable",
  components: { PrincipalAvatar },
  props: {
    repositoryList: {
      required: true,
      type: Object as PropType<Repository[]>,
    },
  },
  setup(props, ctx) {
    const router = useRouter();

    const clickRepository = function (section: number, row: number) {
      const repository = props.repositoryList[row];
      router.push(
        `/project/${projectSlug(repository.project)}#version-control`
      );
    };

    return {
      columnList,
      clickRepository,
    };
  },
};
</script>
