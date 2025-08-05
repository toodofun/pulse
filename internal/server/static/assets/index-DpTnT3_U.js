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

import{k as s,j as l}from"./index-BpEXRqn4.js";const p={list:()=>s.get("/monitor"),delete:e=>s.delete(`/monitor/${e}`),enable:e=>s.put(`/monitor/${e}/enable`),disable:e=>s.put(`/monitor/${e}/disable`)};function i(e){var a,n,t="";if(typeof e=="string"||typeof e=="number")t+=e;else if(typeof e=="object")if(Array.isArray(e)){var r=e.length;for(a=0;a<r;a++)e[a]&&(n=i(e[a]))&&(t&&(t+=" "),t+=n)}else for(n in e)e[n]&&(t&&(t+=" "),t+=n);return t}function o(){for(var e,a,n=0,t="",r=arguments.length;n<r;n++)(e=arguments[n])&&(a=i(e))&&(t&&(t+=" "),t+=a);return t}const u={succeeded:{pingClass:"bg-green-300",dotClass:"bg-green-400",pingHex:"#6ee7b7",dotHex:"#4ade80"},failed:{pingClass:"bg-red-300",dotClass:"bg-red-400",pingHex:"#fca5a5",dotHex:"#f87171"},running:{pingClass:"bg-sky-300",dotClass:"bg-sky-400",pingHex:"#7dd3fc",dotHex:"#38bdf8"},pending:{pingClass:"bg-gray-300",dotClass:"bg-gray-400",pingHex:"#d1d5db",dotHex:"#9ca3af"},unknown:{pingClass:"bg-gray-300",dotClass:"bg-gray-400",pingHex:"#d1d5db",dotHex:"#9ca3af"}},f=({type:e="running",animation:a=!0,size:n=12})=>{const{pingClass:t,dotClass:r,pingHex:g,dotHex:d}=u[e];return l.jsxs("span",{"aria-label":e,title:e,className:"relative flex",style:{minWidth:n,minHeight:n,width:n,height:n},children:[l.jsx("span",{className:o("absolute inline-flex h-full w-full rounded-full opacity-75",a&&"animate-ping",t),style:{backgroundColor:`var(--bg-ping, ${g})`}}),l.jsx("span",{className:o("relative inline-flex rounded-full",r),style:{backgroundColor:`var(--bg-dot, ${d})`,width:n,height:n}})]})};export{f as A,p as m};
