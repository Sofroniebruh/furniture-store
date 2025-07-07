<script setup>

import ProductCard from "@/components/ProductCard.vue";
import {productsData2} from "@/lib/data.js";
import Wrapper from "@/components/Wrapper.vue";
import Sorting from "@/components/Sorting.vue";
import {SlidersHorizontal} from 'lucide-vue-next';
import {Button} from "@/components/ui/button/index.js";
import SheetGeneral from "@/components/SheetGeneral.vue";
</script>

<template>
  <div class="flex flex-col md:flex-row">
    <Wrapper class="md:hidden block">
      <SheetGeneral title="FUMI" side="left">
        <template #trigger>
          <Button class="cursor-pointer" variant="outline">Filters
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
      <div class="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-6 products-block">
        <ProductCard class="w-full max-w-[300px]" v-for="(product, index) in productsData2" :key="index"
                     :product-price="product.productPrice"
                     :product-image="product.productPicture" :product-colors-amount="product.amountOfColors"
                     :product-name="product.productName"/>
      </div>
    </Wrapper>
  </div>
</template>

<style scoped>

</style>