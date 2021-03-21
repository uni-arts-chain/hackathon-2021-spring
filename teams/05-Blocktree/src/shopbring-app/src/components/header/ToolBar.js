import React from "react";
import { Menu, Button, Icon } from "semantic-ui-react";

import PopupLogin from "../popup/PopupLogin";
import PopupCart from "../popup/PopupCart";
import UserInfo from "./UserInfo";

const dataToolBtnsRight = [
  {
    title: "userinfo",
    icon: "user",
  },
  {
    title: "favorite",
    icon: "heart",
  },
  {
    title: "message",
    icon: "bell",
  },
  {
    title: "cart",
    icon: "shop",
  },
  {
    title: "help",
    icon: "question circle outline",
  },
];

// Tool Button 工具栏按钮
const ToolBtn = (index, feature, icon) => (
  <Button className="tool-btn" animated="fade" floated="right" tabIndex={index}>
    <Button.Content as="span" hidden>
      {feature}
    </Button.Content>
    <Button.Content as="span" visible>
      <Icon name={icon} />
    </Button.Content>
  </Button>
);

const ToolBar = (props) => {
  let itemsToolRight = dataToolBtnsRight.map(function (item, index) {
    let title = item.title
      .replace(/\-/g, " ")
      .replace(/^[a-z]?/, function ($0) {
        return $0.toUpperCase();
      });

    let popupWrap = null;

    let btnTool = ToolBtn(index, title, item.icon);

    switch (item.title) {
      case "userinfo":
        popupWrap = <UserInfo />;
        break;
      case "favorite":
        popupWrap = btnTool;
        break;
      case "message":
        popupWrap = btnTool;
        break;
      case "cart":
        popupWrap = <PopupCart cart={props.cart} btnTool={btnTool} />;
        break;
    }
    return (
      <Menu.Item key={index} as="li">
        {popupWrap}
      </Menu.Item>
    );
  });

  return (
    <Menu.Menu as="ul" position="right">
      {itemsToolRight}
    </Menu.Menu>
  );
};

export default ToolBar;
