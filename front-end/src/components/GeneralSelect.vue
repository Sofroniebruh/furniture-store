<script setup lang="ts">
import {
  Select,
  SelectContent,
  SelectGroup,
  SelectItem,
  SelectLabel,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import {PropType} from "vue";

export interface ProductEvent {
  name: string;
}

const props = defineProps({
  selectValue: {
    type: String
  },
  selectLabel: {
    type: String
  },
  events: {
    type: Object as PropType<ProductEvent[]>
  },
  productEvent: {
    type: String
  }
})

const emit = defineEmits(["update:productEvent"])

const handleEvent = (event: string) => {
  emit("update:productEvent", event)
}
</script>

<template>
  <Select>
    <SelectTrigger class="w-full h-full">
      <SelectValue :placeholder="props.selectValue"/>
    </SelectTrigger>
    <SelectContent>
      <SelectGroup>
        <SelectLabel>{{ props.selectLabel }}</SelectLabel>
        <SelectItem @click="handleEvent(event.name)" v-for="(event, index) in props.events"
                    :value="event.name" :key="index">
          {{ event.name }}
        </SelectItem>
      </SelectGroup>
    </SelectContent>
  </Select>
</template>