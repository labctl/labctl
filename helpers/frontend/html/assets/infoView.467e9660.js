import{c as $,a as v,b as l,d as o,e as b,u as y,f as x,g as z,h as B,i as d,F as _,j as V,k as P,l as R,m as h,w as p,n as s,o as m,p as u,q as D,s as g}from"./index.de198c58.js";const N=$("divider",`
 position: relative;
 display: flex;
 width: 100%;
 box-sizing: border-box;
 font-size: 16px;
 color: var(--n-text-color);
 transition:
 color .3s var(--n-bezier),
 background-color .3s var(--n-bezier);
`,[v("vertical",`
 margin-top: 24px;
 margin-bottom: 24px;
 `,[v("no-title",`
 display: flex;
 align-items: center;
 `)]),l("title",`
 display: flex;
 align-items: center;
 margin-left: 12px;
 margin-right: 12px;
 white-space: nowrap;
 font-weight: var(--n-font-weight);
 `),o("title-position-left",[l("line",[o("left",{width:"28px"})])]),o("title-position-right",[l("line",[o("right",{width:"28px"})])]),o("dashed",[l("line",`
 background-color: #0000;
 height: 0px;
 width: 100%;
 border-style: dashed;
 border-width: 1px 0 0;
 `)]),o("vertical",`
 display: inline-block;
 height: 1em;
 margin: 0 8px;
 vertical-align: middle;
 width: 1px;
 `),l("line",`
 border: none;
 transition: background-color .3s var(--n-bezier), border-color .3s var(--n-bezier);
 height: 1px;
 width: 100%;
 margin: 0;
 `),v("dashed",[l("line",{backgroundColor:"var(--n-color)"})]),o("dashed",[l("line",{borderColor:"var(--n-color)"})]),o("vertical",{backgroundColor:"var(--n-color)"})]),T=Object.assign(Object.assign({},x.props),{titlePlacement:{type:String,default:"center"},dashed:Boolean,vertical:Boolean}),f=b({name:"Divider",props:T,setup(t){const{mergedClsPrefixRef:i,inlineThemeDisabled:n}=y(t),a=x("Divider","-divider",N,V,t,i),c=z(()=>{const{common:{cubicBezierEaseInOut:e},self:{color:w,textColor:C,fontWeight:k}}=a.value;return{"--n-bezier":e,"--n-color":w,"--n-text-color":C,"--n-font-weight":k}}),r=n?B("divider",void 0,c,t):void 0;return{mergedClsPrefix:i,cssVars:n?void 0:c,themeClass:r==null?void 0:r.themeClass,onRender:r==null?void 0:r.onRender}},render(){var t;const{$slots:i,titlePlacement:n,vertical:a,dashed:c,cssVars:r,mergedClsPrefix:e}=this;return(t=this.onRender)===null||t===void 0||t.call(this),d("div",{role:"separator",class:[`${e}-divider`,this.themeClass,{[`${e}-divider--vertical`]:a,[`${e}-divider--no-title`]:!i.default,[`${e}-divider--dashed`]:c,[`${e}-divider--title-position-${n}`]:i.default&&n}],style:r},a?null:d("div",{class:`${e}-divider__line ${e}-divider__line--left`}),!a&&i.default?d(_,null,d("div",{class:`${e}-divider__title`},this.$slots),d("div",{class:`${e}-divider__line ${e}-divider__line--right`})):null)}}),j=u("websock state"),E=u("store"),M=m("p",{align:"right"},[u(" Designed 2022 by Johann Kellerman, "),m("a",{href:"//github.com/kellerza"},"@kellerza")],-1),S=b({__name:"infoView",setup(t){const i=P();return(n,a)=>(D(),R(_,null,[h(s(f),null,{default:p(()=>[j]),_:1}),h(s(g),{value:s(i).wsState,copyable:"",boxed:"",color:"",theme:"dark",expanded:"","expand-depth":3},null,8,["value"]),h(s(f),null,{default:p(()=>[E]),_:1}),h(s(g),{value:s(i),copyable:"",boxed:"",color:"",theme:"dark","expand-depth":3,"preview-mode":""},null,8,["value"]),M],64))}});export{S as default};
