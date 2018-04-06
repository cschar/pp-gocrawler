<h1> {{.Title}} </h1>


<h4> image count: {{.ImageCount}} </h4>

<div style="margin:30px; background-color: #e2e0e6">

<h3> Upload file </h3>
<form enctype="multipart/form-data" action="/upload" method="post">
    <input type="file" name="uploadfile" />
    <input type="hidden" name="token" value="{{.}}"/>
    <input type="submit" value="upload" />
</form>
</div>

<a href="/agg"> new agg </a>



