<template>
  <div>
    <h2 class="mx-5 my-5">GoNews - агрегатор новостей.</h2>
    <div v-for="news in news" :key="news.ID">
      <v-card elevation="10" outlined class="mx-5 my-5">
        <v-card-title>
          <a :href="news.Link" target="_blank"> {{ news.Title }} </a>
        </v-card-title>
        <v-card-text>
          {{ news.Content }}
          <v-card-subtitle>
            {{ new Date(news.PubTime * 1000) }}
          </v-card-subtitle>
        </v-card-text>
      </v-card>
    </div>
  </div>
</template>

<script>
export default {
  name: "News",
  data() {
    return {
      news: [],
    };
  },
  mounted() {
    let url = "http://" + window.location.hostname + "/news/40";
    fetch(url)
      .then((response) => response.json())
      .then((data) => (this.news = data));
  },
};
</script>

<style scoped>
</style>
