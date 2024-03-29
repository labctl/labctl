import{c as $,a as v,b as n,d as l,e as b,u as y,f as x,g as z,h as B,i as d,F as _,j as V,k as P,l as R,m as h,w as p,n as s,o as m,p as u,q as D,s as g}from"./index-2265a247.js";const N=$("divider",`
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
 `)]),n("title",`
 display: flex;
 align-items: center;
 margin-left: 12px;
 margin-right: 12px;
 white-space: nowrap;
 font-weight: var(--n-font-weight);
 `),l("title-position-left",[n("line",[l("left",{width:"28px"})])]),l("title-position-right",[n("line",[l("right",{width:"28px"})])]),l("dashed",[n("line",`
 background-color: #0000;
 height: 0px;
 width: 100%;
 border-style: dashed;
 border-width: 1px 0 0;
 `)]),l("vertical",`
 display: inline-block;
 height: 1em;
 margin: 0 8px;
 vertical-align: middle;
 width: 1px;
 `),n("line",`
 border: none;
 transition: background-color .3s var(--n-bezier), border-color .3s var(--n-bezier);
 height: 1px;
 width: 100%;
 margin: 0;
 `),v("dashed",[n("line",{backgroundColor:"var(--n-color)"})]),l("dashed",[n("line",{borderColor:"var(--n-color)"})]),l("vertical",{backgroundColor:"var(--n-color)"})]),T=Object.assign(Object.assign({},x.props),{titlePlacement:{type:String,default:"center"},dashed:Boolean,vertical:Boolean}),f=b({name:"Divider",props:T,setup(t){const{mergedClsPrefixRef:i,inlineThemeDisabled:o}=y(t),a=x("Divider","-divider",N,V,t,i),c=z(()=>{const{common:{cubicBezierEaseInOut:e},self:{color:w,textColor:C,fontWeight:k}}=a.value;return{"--n-bezier":e,"--n-color":w,"--n-text-color":C,"--n-font-weight":k}}),r=o?B("divider",void 0,c,t):void 0;return{mergedClsPrefix:i,cssVars:o?void 0:c,themeClass:r==null?void 0:r.themeClass,onRender:r==null?void 0:r.onRender}},render(){var t;const{$slots:i,titlePlacement:o,vertical:a,dashed:c,cssVars:r,mergedClsPrefix:e}=this;return(t=this.onRender)===null||t===void 0||t.call(this),d("div",{role:"separator",class:[`${e}-divider`,this.themeClass,{[`${e}-divider--vertical`]:a,[`${e}-divider--no-title`]:!i.default,[`${e}-divider--dashed`]:c,[`${e}-divider--title-position-${o}`]:i.default&&o}],style:r},a?null:d("div",{class:`${e}-divider__line ${e}-divider__line--left`}),!a&&i.default?d(_,null,d("div",{class:`${e}-divider__title`},this.$slots),d("div",{class:`${e}-divider__line ${e}-divider__line--right`})):null)}}),j=m("p",{align:"right"},[u(" Designed 2022 by Johann Kellerman, "),m("a",{href:"//github.com/kellerza"},"@kellerza")],-1),M=b({__name:"infoView",setup(t){const i=P();return(o,a)=>(D(),R(_,null,[h(s(f),null,{default:p(()=>[u("websock state")]),_:1}),h(s(g),{value:s(i).wsState,copyable:"",boxed:"",color:"",theme:"dark",expanded:"","expand-depth":3},null,8,["value"]),h(s(f),null,{default:p(()=>[u("store")]),_:1}),h(s(g),{value:s(i),copyable:"",boxed:"",color:"",theme:"dark","expand-depth":3,"preview-mode":""},null,8,["value"]),j],64))}});export{M as default};
