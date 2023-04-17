// this is internal and not intended to be run on its own

import render from 'preact-render-to-string/jsx';
import { h } from 'preact';


const App = <div data-foo={true} />;

console.log(render(App));