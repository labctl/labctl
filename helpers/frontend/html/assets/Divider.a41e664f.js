import{g as p,h as c,i as n,j as r,d as f,k as u,l as v,m as b,n as x,p as s,F as C,q as $}from"./index.3e806caa.js";var w=p("divider",`
 position: relative;
 display: flex;
 width: 100%;
 box-sizing: border-box;
 font-size: 16px;
 color: var(--n-text-color);
 transition:
 color .3s var(--n-bezier),
 background-color .3s var(--n-bezier);
`,[c("vertical",`
 margin-top: 24px;
 margin-bottom: 24px;
 `,[c("no-title",`
 display: flex;
 align-items: center;
 `)]),n("title",`
 display: flex;
 align-items: center;
 margin-left: 12px;
 margin-right: 12px;
 white-space: nowrap;
 font-weight: var(--n-font-weight);
 `),r("title-position-left",[n("line",[r("left",{width:"28px"})])]),r("title-position-right",[n("line",[r("right",{width:"28px"})])]),r("dashed",[n("line",`
 background-color: #0000;
 height: 0px;
 width: 100%;
 border-style: dashed;
 border-width: 1px 0 0;
 `)]),r("vertical",`
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
 `),c("dashed",[n("line",{backgroundColor:"var(--n-color)"})]),r("dashed",[n("line",{borderColor:"var(--n-color)"})]),r("vertical",{backgroundColor:"var(--n-color)"})]);const _=Object.assign(Object.assign({},v.props),{titlePlacement:{type:String,default:"center"},dashed:Boolean,vertical:Boolean});var z=f({name:"Divider",props:_,setup(t){const{mergedClsPrefixRef:l,inlineThemeDisabled:o}=u(t),d=v("Divider","-divider",w,$,t,l),a=b(()=>{const{common:{cubicBezierEaseInOut:e},self:{color:h,textColor:g,fontWeight:m}}=d.value;return{"--n-bezier":e,"--n-color":h,"--n-text-color":g,"--n-font-weight":m}}),i=o?x("divider",void 0,a,t):void 0;return{mergedClsPrefix:l,cssVars:o?void 0:a,themeClass:i==null?void 0:i.themeClass,onRender:i==null?void 0:i.onRender}},render(){var t;const{$slots:l,titlePlacement:o,vertical:d,dashed:a,cssVars:i,mergedClsPrefix:e}=this;return(t=this.onRender)===null||t===void 0||t.call(this),s("div",{role:"separator",class:[`${e}-divider`,this.themeClass,{[`${e}-divider--vertical`]:d,[`${e}-divider--no-title`]:!l.default,[`${e}-divider--dashed`]:a,[`${e}-divider--title-position-${o}`]:l.default&&o}],style:i},d?null:s("div",{class:`${e}-divider__line ${e}-divider__line--left`}),!d&&l.default?s(C,null,s("div",{class:`${e}-divider__title`},this.$slots),s("div",{class:`${e}-divider__line ${e}-divider__line--right`})):null)}});export{z as N};
