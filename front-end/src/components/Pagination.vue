<script setup>

import {Button} from "@/components/ui/button/index.js";
import {ChevronLeft, ChevronRight} from 'lucide-vue-next';
import {useProducts} from "@/composables/useProducts.js";
import {cn} from "@/lib/utils.js";
import {computed} from "vue";

const {setPage, currentPage, totalPages, handleNext, handlePrevious} = useProducts()
const range = computed(() => Array.from({length: totalPages.value}))

const prev = () => {
  window.scrollTo({top: 0})
  handlePrevious(currentPage.value - 1)
}
const next = () => {
  window.scrollTo({top: 0})
  handleNext(currentPage.value + 1)
}
const set = (pageNumber) => {
  window.scrollTo({top: 0})
  setPage(pageNumber)
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
          @click="set(index + 1)"
          :class="cn('p-2 px-4 rounded-sm hover:bg-gray-100 cursor-pointer', index + 1 === currentPage ? 'border' : '')"
          v-for="(page, index) in range"
          :key="index">
        <p>{{ index + 1 }}</p>
      </div>
    </div>
    <Button @click="next" :disabled="currentPage >= totalPages" class="cursor-pointer"
            variant="ghost" size="lg">
      <ChevronRight/>
    </Button>
  </div>
</template>