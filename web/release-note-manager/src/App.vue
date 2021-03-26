<template>
  <v-app>
    <v-app-bar app color="primary" dark>
      <div class="d-flex align-center">
        <v-app-bar-title>Release Note Manager</v-app-bar-title>
      </div>
    </v-app-bar>
    <v-main>
      <ul v-for="repo in repos">
        <li>repo</li>
      </ul>
    </v-main>
  </v-app>
</template>

<script lang="ts">
import { Octokit } from '@octokit/rest';
import Vue from 'vue'
import { Component } from 'vue-property-decorator'
import {RestEndpointMethodTypes} from "@octokit/plugin-rest-endpoint-methods/dist-types/generated/parameters-and-response-types";

const octokit = new Octokit();
octokit.
createOAuthAppAuth({
  clientType: "oauth-app",
  clientId: "fbc0ed44da0ffacc9a7f"
})

@Component({})
export default class App extends Vue {
  private repos!: RestEndpointMethodTypes["repos"]["listForAuthenticatedUser"]["response"];

  public async created() {
    const urlParams = new URLSearchParams(window.location.search)

    urlParams.get()
  }
}
</script>
