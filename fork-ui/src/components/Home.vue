<template>
  <div class="">
    <h1>{{ msg }}</h1>
    <input id="fork-input" v-model="url" placeholder="https://exemple.com"/>
    <br><br>
    <input id="update-input" v-model="updateUrl" placeholder="https://new-url.com" :hidden="hideUpdateInput"/>
    <br>
    <button id="fork-button" @click="generateLink" :disabled="generateBtnDisable" :class="{buttonload: loading}">
      Generate <span :hidden="!loading"><i class="fa fa-refresh fa-spin"></i></span>
    </button>
    <button id="update-button" @click="updateLink" :disabled="updateBtnDisable">Update</button>
    <button id="update-button" @click="infoLink" :disabled="infoBtnDisable">Info</button>
    <br>
    <div :hidden="genUrl == ''">
      <a :href="genUrl" class="text-dark" target="_blank" rel="noopener noreferrer" ref="genUrl">
        {{ genUrl }} 
      </a>
      <button @click="copyURL" id="fork-button"><i class="fa fa-copy"></i></button>
    </div>
    <div id="info-list">
      <ul :hidden="hideInfo">
        <p>URL: {{ linkInfo["url"] }}</p>
        <p>Short: {{ linkInfo["short"] }}</p>
        <p>Click: {{ linkInfo["click"] }}</p>
        <p>UpdateCount : {{ linkInfo["updateCount"] }}</p>
        <p>CreatedAt: {{ linkInfo["createdAt"] }}</p>
        <p>LastViewed: {{ linkInfo["lastViewed"] }}</p>
        <p>LastUpdated: {{ linkInfo["lastUpdated"] }}</p>
      </ul>
    </div>
  </div>
</template>

<script>
export default {
  name: 'Home',
  props: {
    msg: String
  },

  data() {
    let apiURL = "https://api.fork.pw";
    if (window.location.href.includes("localhost")) {
      apiURL = "http://localhost:8080"
    }
      
    return {
      url: "",
      genUrl: "",
      updateUrl: "",
      baseURL: apiURL,
      generateBtnDisable: false,
      updateBtnDisable: true,
      hideUpdateInput: true,
      infoBtnDisable: true,
      hideInfo: true,
      loading: false,
      linkInfo: {
        "url": "",
        "short": "",
        "click": 0,
        "updateCount": 0,
        "createdAt": 0,
        "lastViewed": 0,
        "lastUpdated": 0,
      },
      urlPattern: new RegExp('^(https?:\\/\\/)?'+ // protocol
        '((([a-z\\d]([a-z\\d-]*[a-z\\d])*)\\.)+[a-z]{2,}|'+ // domain name
        '((\\d{1,3}\\.){3}\\d{1,3}))'+ // OR ip (v4) address
        '(\\:\\d+)?(\\/[-a-z\\d%_.~+]*)*'), // fragment locator

      linkPattern: new RegExp(`^${apiURL}\\/[a-zA-Z0-9]{5}`)
    }
  },
  methods: {
    generateLink() {
      if (!this.isURL(this.url)) {
        this.genUrl = "Invalid URL"
        return
      }
      this.loading = true;
      fetch(`${this.baseURL}/gen-link`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ "url": this.url })
      }).then(async (res) => {
          const json = await res.json();
          this.genUrl = json.genURL;
          this.loading = false;
      })
    },
    updateLink() {
      if (!this.isURL(this.updateUrl)) {
        this.genUrl = "Invalid new URL"
        return
      }
      fetch(`${this.baseURL}/update-link`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ "shortUrl": this.url, "newUrl": this.updateUrl })
      }).then(async (res) => {
        res.json().then((data) => {
          this.updateInfo(data);
          this.updateUrl = "";
          this.genUrl = this.url;
        });
      })
    },
    infoLink() {
      let ss = this.url.split("/");
      this.hideInfo = false;
      fetch(`${this.baseURL}/info-link/${ss[ss.length-1]}`, {
        method: "GET",
      }).then(async (res) => {
        res.json().then((data) => {
          this.updateInfo(data);
        });
      })
    },
    updateInfo(json) {
      console.log(json) 
      this.linkInfo = {
        "url": json["url"],
        "short": json["short"],
        "click": json["click"],
        "updateCount": json["updateCount"],
        "createdAt": this.convertUnixTimeStamp(json["createdAt"]),
        "lastViewed": this.convertUnixTimeStamp(json["lastViewed"]),
        "lastUpdated": this.convertUnixTimeStamp(json["lastUpdated"]),
      }
      this.hideInfo = false;
    },
    isURL(str) {
      return !!this.urlPattern.test(str);
    },
    isLink(str) {
      return !!this.linkPattern.test(str)
    },
    convertUnixTimeStamp(unixTime) {
      let d = new Date(unixTime * 1000).toLocaleString("en-US")
      return d;
    },
    async copyURL() {
      await navigator.clipboard.writeText(this.genUrl);
      alert("Copied to clipboard");
    }
  },
  watch: {
    url(val) {
      if (this.isLink(val)) {
        this.generateBtnDisable = true;
        this.infoBtnDisable = false;
        this.hideUpdateInput = false;
        this.genUrl = "";
        this.updateBtnDisable = false;
      } else {
        this.generateBtnDisable = false;
        this.updateBtnDisable = true;
        this.infoBtnDisable = true;
        this.hideUpdateInput = true;
        this.hideInfo = true;
      }
    }
  }

}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
/* h3 {
  margin: 40px 0 0;
}
a {
  color: #42b983;
} */
#info-list {
  text-align: left;
  margin-left: auto;
  margin-right: auto;
  display: table;
}
#fork-button {
    width: auto;
    border-radius: 5px;
    border: 0px;
    padding: 0.5rem 1rem;
    color: #ffffff;
    background: #42b983;
    margin: 1rem;
}
#fork-button:disabled {
    width: auto;
    border-radius: 5px;
    border: 0px;
    padding: 0.5rem 1rem;
    color: #ffffff;
    background: #bfc4d4;
    margin: 1rem;
}
#fork-button:active {
  background-color: #42b983;
  box-shadow: 0 3px #666;
  transform: translateY(2px);
}
#fork-input, #update-input {
    width: 30%;
    min-width: 250px;
    padding: 0.5rem;
    border: 1px solid #ccc;
    border-radius: 5px;
}
#update-button {
    width: auto;
    border-radius: 5px;
    border: 0px;
    padding: 0.5rem 1rem;
    color: #ffffff;
    background: #425eb9;
    margin: 1rem;
}
#update-button:disabled {
    width: auto;
    border-radius: 5px;
    border: 0px;
    padding: 0.5rem 1rem;
    color: #ffffff;
    background: #bfc4d4;
    margin: 1rem;
}
#update-button:active {
  background-color: #425eb9;
  box-shadow: 0 5px #2f51be;
  transform: translateY(4px);
}
pre {outline: 1px solid #ccc; padding: 5px; margin: 5px; }
  .string { color: green; }
  .number { color: darkorange; }
  .boolean { color: blue; }
  .null { color: magenta; }
  .key { color: red; }

.buttonload {
  background-color: #04AA6D; /* Green background */
  border: none; /* Remove borders */
  color: white; /* White text */
  padding: 12px 24px; /* Some padding */
  font-size: 16px; /* Set a font-size */
}
</style>
