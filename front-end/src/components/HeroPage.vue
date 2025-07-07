<script setup>
import Wrapper from "@/components/Wrapper.vue";
import ReviewCard from "@/components/ReviewCard.vue";
import {productsData, reviews} from "@/lib/data.js";
import ProductCard from "@/components/ProductCard.vue";
import {Button} from "@/components/ui/button/index.js";
import {ArrowUp} from 'lucide-vue-next';


const scrollToTop = () => {
  window.scrollTo({
    top: 0,
    behavior: "smooth",
  });
};
</script>

<template>
  <section id="hero-image">
    <div class="relative flex justify-center items-center">
      <div class="relative w-full sm:h-[600px] h-[400px]">
        <img class="w-full object-cover object-center h-full" src="/images/hero-image-main-page.webp"
             alt="hero about image">
      </div>
      <h1 class="absolute hero-text font-semibold text-white text-6xl sm:text-8xl tracking-tight text-balance !leading-tight">
        Furniture for <span
          class="hero-text tracking-tight relative font-bold text-6xl sm:text-8xl whitespace-nowrap bg-white text-[#c9a275] px-4">every kind</span>
        of taste
      </h1>
    </div>
  </section>
  <section id="quality">
    <Wrapper class="mt-[34px] md:mt-[68px] py-0">
      <h1 class="mb-[34px] md:mb-[68px] text-center sm:text-4xl text-3xl text-gray-800 block xl:hidden">Exceptional
        quality.</h1>
      <div class="flex flex-col-reverse xl:flex-row">
        <div class="flex xl:w-1/2 w-full flex-col justify-center items-start gap-6 mt-[34px] xl:mt-0">
          <h1 class="text-center sm:text-4xl text-3xl text-gray-800 xl:block hidden">Exceptional quality.</h1>
          <p class="text-balance text-gray-600 font-medium">At <span
              class="font-semibold text-gray-600">FUMI</span>, we take pride in crafting furniture that stands the
            test of time. Each piece is thoughtfully designed and meticulously constructed using premium materials and
            expert craftsmanship. Whether it’s the smooth finish of a dining table or the plush comfort of a sofa, every
            detail is a testament to our commitment to quality.
          </p>
        </div>
        <div class="xl:w-1/2 w-full">
          <div class="xl:h-[600px] h-[400px] w-full overflow-hidden rounded-sm">
            <img class="object-cover h-full w-full" src="/images/quality-section-image-about-page.webp"
                 alt="quality section image"/>
          </div>
        </div>
      </div>
    </Wrapper>
  </section>
  <section id="reviews">
    <Wrapper class="mt-[34px] md:mt-[68px] py-0">
      <h1 class="mb-[34px] md:mb-[68px] text-center sm:text-4xl text-3xl text-gray-800 block">What people
        say</h1>
      <div class="flex flex-col gap-5 sm:gap-7 justify-center items-center">
        <ReviewCard v-for="(review, index) in reviews" :username="review.username" :rating="review.rating"
                    :review="review.review" :key="index" class="w-full sm:max-w-xl"/>
      </div>
    </Wrapper>
  </section>
  <section id="featured-products">
    <Wrapper class="mt-[34px] md:mt-[68px] flex flex-col justify-center py-0 gap-5 sm:gap-7">
      <h1 class="mb-[34px] md:mb-[68px] text-center sm:text-4xl text-3xl text-gray-800 block">Take a look</h1>
      <div
          class="flex items-center justify-center flex-col sm:grid sm:grid-cols-2 md:grid-cols-3 gap-6 sm:gap-8">
        <div v-for="(product, index) in productsData" class="w-full flex items-center justify-center">
          <ProductCard :key="index" :product-name="product.productName"
                       :product-colors-amount="product.amountOfColors" :product-image="product.productPicture"
                       :product-price="product.productPrice"></ProductCard>
        </div>
      </div>
    </Wrapper>
  </section>
  <section id="to-the-top-button">
    <div class="w-full flex items-center justify-center mt-[34px] sm:mt-[68px]">
      <Button variant="outline" class="hover:bg-gray-100 cursor-pointer text-base text-center" @click="scrollToTop">Back to the top
        <ArrowUp/>
      </Button>
    </div>
  </section>
</template>