<script setup>

import ProductCard from "@/components/ProductCard.vue";
import Wrapper from "@/components/Wrapper.vue";
import Sorting from "@/components/Sorting.vue";
import {SlidersHorizontal} from 'lucide-vue-next';
import {Button} from "@/components/ui/button/index.js";
import SheetGeneral from "@/components/SheetGeneral.vue";
import {useScreenSheetStore} from "@/stores/useScreenSheetStore.ts";
import {useProducts} from "@/composables/useProducts.js";
import {onMounted, watch} from "vue";
import Pagination from "@/components/Pagination.vue";

const smallScreenSheet = useScreenSheetStore()
const {products, loading, error, fetchProducts, totalPages} = useProducts()

onMounted(async () => {
  await fetchProducts()
})

watch(products.value, (newValue) => {
  products.value = newValue;
}, {immediate: true})

</script>

<template>
  <div class="flex flex-col md:flex-row">
    <Wrapper class="md:hidden block">
      <SheetGeneral :is-open="smallScreenSheet.isSheetOpen('SmallScreenFiltering')" title="FUMI" side="left">
        <template #trigger>
          <Button @click="smallScreenSheet.setOpenSheet(true, {name: 'SmallScreenFiltering'})" class="cursor-pointer"
                  variant="outline">Filters
            <SlidersHorizontal/>
          </Button>
        </template>
        <template #content>
          <Wrapper class="pt-0">
            <Sorting/>
          </Wrapper>
        </template>
      </SheetGeneral>
    </Wrapper>
    <Wrapper class="hidden w-1/4 md:block pr-0">
      <Sorting/>
    </Wrapper>
    <Wrapper class="p-0 md:pl-5 md:pt-5 pb-0 w-full md:w-3/4">
      <div v-if="loading" class="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-6 products-block">
        <div v-for="n in 8" :key="n" class="flex w-full items-center justify-center p-5 md:p-0">
          <div class="w-full max-w-[300px] h-64 bg-gray-200 animate-pulse rounded"></div>
        </div>
      </div>
      <div class="flex justify-center items-center md:h-full h-[calc(100vh-200px-68px-76px)]"
           v-else-if="products === undefined || products === null">
        <h1 class="font-semibold text-base sm:text-lg">No products matching your request</h1>
      </div>
      <div v-else-if="error" class="text-center py-8">
        <p class="text-red-600">Failed to load products. Please try again.</p>
        <Button @click="fetchProducts" class="mt-4">Retry</Button>
      </div>
      <div class="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-6 products-block">
        <div v-for="(product, index) in products" :key="index"
             class="flex w-full items-center justify-center p-5 md:p-0">
          <ProductCard class="w-full max-w-[300px]"
                       :product-price="product.price"
                       :product-image="product.pictureUrls[0]" :product-colors-amount="product.colors.length"
                       :product-name="product.name"/>
        </div>
      </div>
      <Pagination v-if="products !== null"/>
    </Wrapper>
  </div>
</template>