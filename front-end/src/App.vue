<script setup>
import Header from './components/Header.vue'
import Footer from "@/components/Footer.vue";
</script>

<template>
  <div class="w-full">
    <Header/>
    <div class="mt-[68px]">
      <router-view></router-view>
    </div>
    <Footer/>
  </div>
</template>
