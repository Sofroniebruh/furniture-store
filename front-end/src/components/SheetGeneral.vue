<script setup>
import {Sheet, SheetContent, SheetDescription, SheetHeader, SheetTitle, SheetTrigger,} from '@/components/ui/sheet'

const props = defineProps({
  title: {
    type: String,
    required: false,
  },
  isOpen: {
    type: Boolean,
  },
  description: {
    type: String,
    required: false,
  },
  side: {
    type: String,
    default: 'right',
    validator: (value) => ["left", "right", "bottom", "top"].includes(value),
  },
  triggerClass: {
    type: String,
    required: false,
  }
})
</script>

<template>
  <Sheet :open="isOpen">
    <SheetTrigger :class="triggerClass">
      <slot name="trigger"/>
    </SheetTrigger>
    <SheetContent :side="side">
      <SheetHeader>
        <SheetTitle>{{ title }}</SheetTitle>
        <SheetDescription>
          {{ description }}
        </SheetDescription>
      </SheetHeader>
      <slot name="content"/>
    </SheetContent>
  </Sheet>
</template>