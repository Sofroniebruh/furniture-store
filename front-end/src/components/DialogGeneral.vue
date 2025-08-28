<script setup lang="ts">
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from '@/components/ui/dialog'

const props = defineProps({
  title: {
    type: String,
    required: false,
  },
  description: {
    type: String,
    required: false,
  },
  footer: {
    type: Boolean,
    default: false,
  },
  isOpen: {
    type: Boolean,
    required: true,
  }
})
</script>

<template>
  <Dialog :open="props.isOpen">
    <DialogTrigger class="w-full">
      <div class="flex">
        <slot name="trigger"/>
      </div>
    </DialogTrigger>
    <DialogContent>
      <DialogHeader>
        <DialogTitle>{{ props.title }}</DialogTitle>
        <DialogDescription>
          {{ props.description }}
        </DialogDescription>
      </DialogHeader>
      <slot name="content"/>
      <DialogFooter v-if="props.footer">
        <slot name="footer"/>
      </DialogFooter>
    </DialogContent>
  </Dialog>
</template>