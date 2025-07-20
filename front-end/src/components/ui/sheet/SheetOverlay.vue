<script setup>
import {reactiveOmit} from "@vueuse/core";
import {DialogOverlay} from "reka-ui";
import {cn} from "@/lib/utils";
import {useScreenSheetStore} from "@/stores/useScreenSheetStore.ts";

const props = defineProps({
  forceMount: {type: Boolean, required: false},
  asChild: {type: Boolean, required: false},
  as: {type: null, required: false},
  class: {type: null, required: false},
});

const delegatedProps = reactiveOmit(props, "class");
const smallScreenSheet = useScreenSheetStore()
</script>

<template>
  <DialogOverlay
      data-slot="sheet-overlay"
      @click="smallScreenSheet.setAllSheetsClosed()"
      :class="
      cn(
        'data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0 fixed inset-0 z-50 bg-gray-400/40',
        props.class,
      )
    "
      v-bind="delegatedProps"
  >
    <slot/>
  </DialogOverlay>
</template>
