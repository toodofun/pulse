/**
 * Copyright 2025 The Toodofun Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http:www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

import{r as c,I as i,aF as f,h as l,j as r}from"./index-BpEXRqn4.js";function o(){return o=Object.assign?Object.assign.bind():function(e){for(var t=1;t<arguments.length;t++){var s=arguments[t];for(var n in s)Object.prototype.hasOwnProperty.call(s,n)&&(e[n]=s[n])}return e},o.apply(this,arguments)}const u=(e,t)=>c.createElement(i,o({},e,{ref:t,icon:f})),p=c.forwardRef(u),m=({title:e="Go Back",onClick:t})=>{if(!t){const s=l();t=()=>{s(-1)}}return r.jsx("div",{children:r.jsxs("div",{className:"bg-slate-800 px-4 py-3 flex w-fit items-center rounded-lg gap-2 cursor-pointer active:bg-slate-800 hover:bg-slate-700 select-none",onClick:t,children:[r.jsx(p,{className:"text-[0.5rem]"}),r.jsx("div",{className:"text-sm",children:e})]})})};function d(e){const t=Math.floor(e/3600),s=Math.floor(e%3600/60),n=e%60,a=[];return t>0&&a.push(`${t} hour`),s>0&&a.push(`${s} minute`),(n>0||a.length===0)&&a.push(`${n} second`),a.join(" ")}export{m as G,d as f};
