package doc

const HTML = `
<html>
<body>
<head>
	<style type="text/css">
/* Just some base styles not needed for example to function */
*, html { font-family: Verdana, Arial, Helvetica, sans-serif; }

body, form, ul, li, p, h1, h2, h3, h4, h5
{
	margin: 0;
	padding: 0;
}
body { margin: 0; }
img { border: none; }
p
{
	font-size: 1em;
	margin: 0 0 1em 0;
}

html { font-size: 100%; /* IE hack */ }
body { font-size: 1em; /* Sets base font size to 16px */ }
table { font-size: 100%; /* IE hack */ }
input, select, textarea, th, td { font-size: 1em; }
	
/* CSS Tree menu styles */
ol.tree
{
	padding: 0 0 0 30px;
	width: 300px;
}
	li 
	{ 
		position: relative; 
		margin-left: -15px;
		list-style: none;
	}
	li.file
	{
		margin-left: -1px !important;
	}
		li.file a
		{
			background: url(data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAABAAAAAQCAYAAAAf8/9hAAAAGXRFWHRTb2Z0d2FyZQBBZG9iZSBJbWFnZVJlYWR5ccllPAAAAYdJREFUeNqMksFqwkAQhmc3W4q0OSq09170LdJzb32IvkQVpIeeS6Gehdz6AkUoePOYe71VEBGiVkUxyaYzY9bENIYuTDbZ3f/b+Wci7ptNoCGEuMPpCsrHp9Z6qKMI4jjmBWVecL7utlqdHW4GFFpDhHtm3/M86PZ6D3QXxpchSqIlITSJgoBjt93CZr2G1WoFy+WS5/d2u4PnblF3YwAqQnEyBKElPs8sK006DCFEsO/78Oq64NTrby+u+4g7TwzQKUCSWDIG/kAcx2E7tm3Dh+fNDhmE2QwygCIIZcKvaDe1gBt5C9mRh9AlqCkE7C1ICZDcVASx9gCZWsgADhmUQCycwxMArgHXYZ9PIcTC9fDIAvY8nwHJTkHYQhCcsEAZYKvKIIJrmbGwGI/holo96oIugdC6sbCeTkEO+33wRyP6bVVSCD5kYOabumPW6SxpSEsXVL4Hg9piNjvPVPMgoKJRKCygUoq7QGdJQ1q6tYJR+5lM5peNxjP8Y+jNZk4aKsevAAMAmFzedjV8x2YAAAAASUVORK5CYII=) 0 0 no-repeat;
			padding-left: 21px;
			text-decoration: none;
			display: block;
		}
	li input
	{
		position: absolute;
		left: 0;
		margin-left: 0;
		opacity: 0;
		z-index: 2;
		cursor: pointer;
		height: 1em;
		width: 1em;
		top: 0;
	}
		li input + ol
		{
			background: url(data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAABAAAAAQCAYAAAAf8/9hAAAAGXRFWHRTb2Z0d2FyZQBBZG9iZSBJbWFnZVJlYWR5ccllPAAAAURJREFUeNpi/P//PwMlgImBQkCxASwwRlLLKwYmJqZgRkbGbiBXEYva+0Dvlv7792/tvBoxTAO+fv0MororE6UU9VU5MHRfvP1DsX3+M5DhaxkYxDC98ObNGxBW1FVmY/j16xcYu6SdYvjw4QPDixcvGGSEvoLlQeqweuHdu7dg+vfv32D85ctXsNijR4/B4hwcnHA1WA348uUbmP779y+DUchOuIKQsltgetsUE7garAb8/w9h/vz5h+H0Sk8w2yRsN8OZVa5g9ocPn+BqsBrAzs4PdQEzw48ff+Fi375B2Gxs3HA1WNPB45NlDNzcIvfPXv8LVMwJxmdWOcDZF2//A8uD1GF1wefXZ8Q+Pt42oWN+VBED41d5DKv+/30IlJ8IVCcF5D2DCTPC8gIwAXEDKT4Qk0Di+wzU8xnDgKGbmQACDAAtTZadqmiADQAAAABJRU5ErkJggg==) 40px 0 no-repeat;
			margin: -0.938em 0 0 -44px; /* 15px */
			height: 1em;
		}
		li input + ol > li { display: none; margin-left: -14px !important; padding-left: 1px; }
	li label
	{
		background: url(data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAABAAAAAQCAYAAAAf8/9hAAAAGXRFWHRTb2Z0d2FyZQBBZG9iZSBJbWFnZVJlYWR5ccllPAAAATNJREFUeNqkk71KA1EQhc/dOxsbEYukVYm9jQg+hz6CYGkrBNsEKwtrX0EfwU4UVFDLFWSDlYGAGszP3p91ZlNpdoVrBg572bnf2ZlhR+V5jnmCbo9VZTJS2ODHbkX63Od4Ij7ssdbKbvAFbB76o2GWYTAe42sywX7aQL8/xNnjYsRXttR1G+3tg4tW7twPWGmNXvJcnBvrzRlzyd+c7nTIeWjPXxD9jjqDb0mC9O6+tAdhpwbWQDTTwsCgvrpSCksFr1dsYAsDC1diUETFe11bgC0qcCBnMoiCQikIS9ZxBcawgQnkIwgrBiT9h1YQRYUBkfFTA2/DKvA8RGGJ5xf/OcSK0JogLH2MsCS/VKzjsBZYwtJDiu7nSevyP4v00kNXNqnGWpbFCuQt613Nu87fAgwAb3KTD1NdyNYAAAAASUVORK5CYII=) 15px 1px no-repeat;
		cursor: pointer;
		display: block;
		padding-left: 37px;
	}

	li input:checked + ol
	{
		background: url(data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAABAAAAAQCAYAAAAf8/9hAAAAGXRFWHRTb2Z0d2FyZQBBZG9iZSBJbWFnZVJlYWR5ccllPAAAASxJREFUeNpi/P//PwMlgImBQkCxASwwRlLLKwYmJqZgRkbGbiBXEYva+0Dvlv7792/tvBoxTAO+fv0MororE6UU9VU5MHRfvP1DsX3+M5DhaxkYsBjw5s0bEKWoq6zA8OvXL7AYKIC/f//O8OPHDwYZIVaQGqjLlDENePfuLZj+/fs3GH/58pXh/fv3YDYIcHBwwtVgDYMvX76B6b9//zIYhezEULhtiglcDVYD/v+HMH/+/MNweqUnhsIPHz7B1WA1gJ2dH+oCZqCf/2IoZGPjhqvBmg4enyxj4OYWuX/2+l+gYk4MfPH2P7A8SB1WF3x+fUbs4+NtEzrmRxUxMH6Vx7Dq/9+HQPmJQHVSQN4zmDAjLC8AExA3kOIDMQkkvs9APZ8xDBi6mQkgwADDMYZH9Ls66AAAAABJRU5ErkJggg==) 40px 5px no-repeat;
		margin: -1.25em 0 0 -44px; /* 20px */
		padding: 1.563em 0 0 80px;
		height: auto;
	}
		li input:checked + ol > li { display: block; margin: 0 0 0.125em;  /* 2px */}
		li input:checked + ol > li:last-child { margin: 0 0 0.063em; /* 1px */ }

	div.fd {position:fixed; top:0;_position:relative; z-index:100; _top:expression(offsetParent.scrollTop+0);} 
	</style>
</head>
<div id="main" style="width:1000px;margin:0px;padding:0px">
<div style="width:300px;float:left;height:100%;overflow:auto;" class="fd">
{{range $hkey,$hval:=.Items}}
	<div style="margin-left:10px;"><a href="#{{$hkey}}">{{$hkey}}</a></div>
{{end}}
	<ol class="tree">
		{{.Tree}}
	</ol>
	<br/><br/><br/><br/>
</div>
<div id="content" style="margin-left:300px;float:left;width:800px;">

{{range $hkey,$hval:=.Items}}
<h1 id="{{$hkey}}" style="background:#AABBCC;font-size:25px;">{{$hkey}}</h1>
<div style="margin-left:10px;margin-right:10px;">{{$hval}}</div>
<br/>
{{end}}
{{range $hkey,$hval:=.Apis}}
<h1 id="{{$hkey}}" style="background:#AABBCC;font-size:25px;">{{$hkey}}</h1>
{{range $hval}}
<div style="margin-left:10px;margin-right:10px;">
	<h2 id="{{.Abs}}" style="background:#E0EBF5;font-size:20px;">
		{{if .Marked}}
			{{.Doc.Title}}
		{{else}}
			{{.Name}}
		{{end}}
	</h2>
	<div style="margin-left:10px;margin-right:10px;">
			<div style="margin:0px;padding-left:5px;padding-bottom:10px;background:#EEEBF5;">
				<div style="padding-top:10px;">Path: <a href="#{{.Abs}}">{{.Pkg}}/{{.Name}}</a></div>
				<div style="margin-top:10px;">Pattern: {{.Pattern}}</div>
				{{if .Marked}}
				{{if .Doc.Url}}
				<div style="margin-top:10px;">Example: 
					{{range .Doc.Url}}
					&nbsp;<a href="{{.}}">{{.}}</a>
					{{end}}
				</div>
				{{end}}
				{{end}}
			</div>
		{{if .Marked}}
			{{if .Doc.Detail}}
			<div style="margin-top:5px;padding-left:5px;padding-top:10px;padding-bottom:10px;background:#F0F0F0;">{{.Doc.Detail}}</div>
			{{end}}
			<br/>

			<b>Parameters(Required)</b>
			<ul style="margin:0px;padding-left:30px;background:#EEE;">
				{{range $key,$val:=.Doc.ArgsR}}
				<li><b style="font-size:16px;">{{$key}}</b> {{$val}}</li>
				{{end}}
			</ul>
			<br/>

			<b>Parameters(Optioned)</b>
			<ul style="margin:0px;padding-left:30px;background:#EEE;">
				{{range $key,$val:=.Doc.ArgsO}}
				<li><b style="font-size:16px;">{{$key}}</b> {{$val}}</li>
				{{end}}
			</ul>
			<br/>

			<b>Parameter Value Option</b>
			<ul style="margin:0px;padding-left:30px;background:#EEE;">
				{{range $key,$val:=.Doc.Option}}
				<li>
					<b style="font-size:16px;">{{$key}}</b>
					<ul style="margin:0px;padding-left:10px;">
						{{range $key1,$val1:=$val}}
						<li><b style="font-size:14px;">{{$key1}}</b> {{$val1}}</li>
						{{end}}
					</ul>
				</li>
				{{end}}
			</ul>
			<br/>

			<b>Return</b>
			<div style="margin:0px;padding:10px;background:#EEE;">
			{{.Doc.ResHTML}}
			</div>

			{{if .Doc.See}}
			<b>See</b>
			<ul style="margin:0px;padding-left:30px;background:#EEE;">
				{{range .Doc.See}}
				<li>
					<a href="#{{.Abs}}">{{.Pkg}}/{{.Name}}</a>
				</li>
				{{end}}
			</ul>
			{{end}}
		{{else}}
			<div style="margin-top:5px;padding-left:5px;padding-top:10px;padding-bottom:10px;background:#F0F0F0;">Not Marked</div>
		{{end}}
	</div>
</div>
<br/><br/>
{{end}}
{{end}}
</div>
</div>
</body>
</html>
`
