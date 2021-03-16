import React from 'react';
import { Grid } from 'semantic-ui-react';
import FooterSubMenu from './FooterSubMenu';
import SubscribeMenu from './SubscribeMenu';

// Footer Submenu Text Data 底部子菜单显示文本的数据
const dataFooterMenu = [
  {
    menuHeader: 'shopbring',
    menuItems: [
      {
        link: 'home',
        path: '/'
      }
    ]
  },
  {
    menuHeader: 'info',
    menuItems: [
      {
        link: 'about-Shopbring',
        path: '/'
      },
      {
        link: 'polkadot',
        path: 'https://polkadot.network/'
      },
    ]
  },
  {
    menuHeader: 'help',
    menuItems: [
      {
        link: 'customer-service',
        path: '/'
      },
      {
        link: 'my-Shopbring',
        path: ''
      },
      {
        link: 'contact',
        path: '/'
      }
    ]
  }
]

const FootMenu = (props) => {
  let itemsSubMenu = dataFooterMenu.map(function(item,index){
    return(
      <Grid.Column key={index}>
        <FooterSubMenu {...item}/>
      </Grid.Column>
    )
  })

  return(
    <div className="footer-menu">
        <Grid columns='four' divided>
          <Grid.Row>
            {itemsSubMenu}
            <SubscribeMenu/>
          </Grid.Row>
        </Grid>
    </div>
  )
}

export default FootMenu;