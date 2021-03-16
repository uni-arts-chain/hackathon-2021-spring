import React, { Component } from "react";
import { Menu, Search } from "semantic-ui-react";
import axios from "axios";

class SearchBar extends Component {
  constructor(props) {
    super(props);
    this.state = {
      isLoading: false,
      loadData: [],
      results: [],
      value: "",
    };
    this.handleResultSelect = this.handleResultSelect.bind(this);
    this.handleSearchChange = this.handleSearchChange.bind(this);
    this.handleSearching = this.handleSearching.bind(this);
    this.handleReachClear = this.handleReachClear.bind(this);
  }

  componentDidMount() {
    // axios
    // .get('/src/data/productData/database.json')
    // .then(res=>this.setState({loadData: res.data}))
    // .catch(err=>console.log(err))
  }

  // 点击搜索结果现在在文本框内
  handleResultSelect(e) {
    this.setState({ value: e.target.innerHTML });
  }

  handleSearchChange(e) {
    this.setState({ isLoading: true, value: e.target.value });

    setTimeout(() => {
      // 文本框内容为空，返回
      if (this.state.value.length < 1) return;

      const regSearch = new RegExp(this.state.value, "gim");

      let resultMatch = new Array();

      this.state.loadData.forEach(function (item, index) {
        if (item.name.search(regSearch) !== -1)
          resultMatch.push({ title: item.name });
      });

      resultMatch =
        resultMatch.length <= 5 ? resultMatch : resultMatch.slice(0, 5);

      this.setState({
        isLoading: false,
        results: resultMatch,
      });
    }, 500);
  }

  handleSearching(e) {
    if (e.keyCode !== 13) return;
    if (e.target.value.trim() === "") return;

    const value = e.target.value;
    const search_url = new URL(value);
    const id = search_url.searchParams.get("id");
    const url = search_url.origin + search_url.pathname + "?id=" + id;

    window.location.href = "/products?url=" + encodeURIComponent(url);
  }

  handleReachClear(e) {
    this.setState({
      isLoading: false,
      results: [],
      value: "",
    });
  }

  render() {
    const { value, results, isLoading } = this.state;

    const { handleResultSelect, handleSearchChange, handleSearching } = this;

    return (
      <Menu.Menu>
        <Search
          className="search-bar"
          placeholder="Search"
          loading={isLoading}
          onResultSelect={handleResultSelect}
          onSearchChange={handleSearchChange}
          onKeyDown={handleSearching}
          results={results}
        />
      </Menu.Menu>
    );
  }
}

export default SearchBar;
