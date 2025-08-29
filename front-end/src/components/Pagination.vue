<script setup>

import {Button} from "@/components/ui/button/index.js";
import {ChevronLeft, ChevronRight} from 'lucide-vue-next';
import {cn} from "@/lib/utils.js";
import {computed} from "vue";

const props = defineProps({
  currentPage: {
    type: Number,
  },
  totalPages: {
    type: Number,
  },
  setPage: {
    type: Function,
  },
  handleNext: {
    type: Function,
  },
  handlePrevious: {
    type: Function,
  }
})

const range = computed(() => Array.from({length: props.totalPages}, (_, i) => i + 1))

const prev = () => {
  window.scrollTo({top: 0, behavior: 'smooth'})
  props.handlePrevious()
}

const next = () => {
  window.scrollTo({top: 0, behavior: 'smooth'})
  props.handleNext()
}

const set = (pageNumber) => {
  window.scrollTo({top: 0, behavior: 'smooth'})
  props.setPage(pageNumber)
}
</script>

<template>
  <div class="flex gap-2 items-center justify-center">
    <Button @click="prev" :disabled="currentPage <= 1" class="cursor-pointer" variant="ghost"
            size="lg">
      <ChevronLeft/>
    </Button>
    <div class="flex gap-1">
      <div
          @click="set(page)"
          :class="cn('p-2 px-4 rounded-sm hover:bg-gray-100 cursor-pointer', page === currentPage ? 'border' : '')"
          v-for="page in range"
          :key="page">
        <p>{{ page }}</p>
      </div>
    </div>
    <Button @click="next" :disabled="currentPage >= totalPages" class="cursor-pointer"
            variant="ghost" size="lg">
      <ChevronRight/>
    </Button>
  </div>
</template>