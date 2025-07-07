<script setup>
import Wrapper from "@/components/Wrapper.vue";
import {ArrowUp, Box, CalendarDays, Truck} from 'lucide-vue-next';
import {Button} from "@/components/ui/button/index.js";

const scrollToTop = () => {
  window.scrollTo({
    top: 0,
    behavior: "smooth",
  });
};
</script>

<template>
  <section id="hero-about-section">
    <div class="relative overflow-hidden flex items-center justify-center">
      <div class="relative w-full sm:h-[600px] h-[400px]">
        <img class="w-full object-cover object-center h-full" src="/images/hero-image-about-page.webp"
             alt="hero image"/>
      </div>
      <div class="bg-black/35 absolute">
        <h1 class="p-0 sm:text-9xl text-7xl text-white">FUMI</h1>
      </div>
    </div>
  </section>
  <section id="statistics"
           class="my-[34px] sm:my-[68px] flex flex-col lg:flex-row justify-around gap-12 lg:gap-0 lg:h-[200px] items-center">
    <div class="flex flex-col items-center">
      <h1 class="text-5xl sm:text-7xl font-extralight">97%</h1>
      <p class="font-medium text-base">customers come back</p>
    </div>
    <div class="flex flex-col items-center">
      <h1 class="text-5xl sm:text-7xl font-extralight">9000+</h1>
      <p class="font-medium text-base">happy customers</p>
    </div>
    <div class="flex flex-col items-center">
      <h1 class="text-5xl sm:text-7xl font-extralight">12000+</h1>
      <p class="font-medium text-base">orders last month</p>
    </div>
  </section>
  <section id="about-info">
    <h1 class="md:hidden text-center block mb-[34px] sm:mb-[64px] text-3xl text-gray-800 w-full">Our passion is
      quality</h1>
    <Wrapper class="flex gap-6 flex-col-reverse md:flex-row py-0">
      <div class="w-full md:w-1/2 flex flex-col justify-center lg:pr-6">
        <div class="flex flex-col gap-6 items-start">
          <h1 class="md:block hidden text-start sm:text-4xl text-gray-800 w-full">Our passion is quality</h1>
          <p class="text-gray-500 text-balance lg:w-[500px]">We have more than 30 years of experience
            producing and exporting
            furniture for international markets. More than <span class="font-semibold text-gray-700">1.000.000 </span>
            pieces of furniture were delivered to families and businesses around the world.
          </p>
        </div>
      </div>
      <div class="w-full md:w-1/2 relative lg:h-[600px] h-[400px] overflow-hidden rounded-sm">
        <img class="object-cover object-center w-full h-full" src="/images/quality-section-image-main-page.webp"
             alt="about image"/>
      </div>
    </Wrapper>
  </section>
  <section id="cards">
    <h1 class="text-center my-[34px] sm:my-[68px] sm:text-4xl text-3xl text-gray-800 w-full">What we offer</h1>
    <Wrapper
        class="py-0 flex justify-center items-center mt-[34px] sm:mt-[68px] w-full lg:justify-around lg:flex-row flex-col lg:gap-0 gap-12">
      <div
          class="flex flex-col rounded-2xl justify-center items-center bg-gray-100/75 sm:p-10 p-7 gap-5 w-full sm:max-w-[300px]">
        <Box class="lg:w-20 lg:h-20 md:w-15 md:h-15 w-12 h-12"></Box>
        <div class="flex justify-center flex-col items-center text-center">
          <h1 class="text-lg font-semibold">Made to Order</h1>
          <p class="text-balance">All pieces made to order for you</p>
        </div>
      </div>
      <div
          class="flex flex-col rounded-2xl justify-center items-center bg-gray-100/75 sm:p-10 p-7 gap-5 w-full sm:max-w-[300px]">
        <Truck class="lg:w-20 lg:h-20 md:w-15 md:h-15 w-12 h-12"></Truck>
        <div class="flex justify-center flex-col items-center text-center">
          <h1 class="text-lg font-semibold">Free Delivery</h1>
          <p class="text-balance">Free delivery for orders world-wide</p>
        </div>
      </div>
      <div
          class="flex flex-col rounded-2xl justify-center items-center bg-gray-100/75 sm:p-10 p-7 gap-5 w-full sm:max-w-[300px]">
        <CalendarDays class="lg:w-20 lg:h-20 md:w-15 md:h-15 w-12 h-12"></CalendarDays>
        <div class="flex justify-center flex-col items-center text-center">
          <h1 class="text-lg font-semibold">Free Exchange</h1>
          <p class="text-balance">Free exchange on all products <br/> during the first week</p>
        </div>
      </div>
    </Wrapper>
  </section>
  <div class="w-full flex items-center justify-center mt-[34px] sm:mt-[68px]">
    <Button variant="outline" class="hover:bg-gray-100 cursor-pointer text-base text-center" @click="scrollToTop">Back
      to the top
      <ArrowUp/>
    </Button>
  </div>
</template>