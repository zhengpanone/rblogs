import os
import sys
import sphinx_rtd_theme
# import os
# import sys
# sys.path.insert(0, os.path.abspath('.'))


# -- Project information -----------------------------------------------------

project = 'GO BLOG'
copyright = '2019, zhengpanone'
author = 'zhengpanone'

# The full version, including alpha/beta/rc tags
release = '1.0'


# -- General configuration ---------------------------------------------------

# Add any Sphinx extension module names here, as strings. They can be
# extensions coming with Sphinx (named 'sphinx.ext.*') or your custom
# ones.
# 'chinese_search', sphinxdrawio.drawio_html

simplepdf_vars = {
    'primary': '#333333',
    'links': '#FF3333',
}

extensions = ['recommonmark', 
              'sphinx.ext.autodoc',
              'sphinx_copybutton',
              'sphinx.ext.napoleon',
              'sphinx.ext.viewcode',
              'sphinx_markdown_tables',
              'sphinx.ext.autosectionlabel',
              'sphinxcontrib.inkscapeconverter',
            #   'sphinx_simplepdf'
              ]


# LaTeX配置
latex_engine = 'xelatex'  # 或者 'pdflatex'，根据你的需求选择
latex_elements = {
    'papersize': 'a4paper',
    'pointsize': '10pt',
    'preamble': '',
    'figure_align': 'htbp',
}

# Add any paths that contain templates here, relative to this directory.
templates_path = ['_templates']

source_suffix = ['.rst', '.md']
#
# This is also used if you do content translation via gettext catalogs.
# Usually you set "language" from the command line for these cases.
language = 'zh_CN'

# List of patterns, relative to source directory, that match files and
# directories to ignore when looking for source files.
# This pattern also affects html_static_path and html_extra_path.
exclude_patterns = []

pygments_style = 'sphinx'

# -- Options for HTML output -------------------------------------------------

# The theme to use for HTML and HTML Help pages.  See the documentation for
# a list of builtin themes.
#
html_theme = 'sphinx_rtd_theme'
html_theme_path = [sphinx_rtd_theme.get_html_theme_path()]
# Add any paths that contain custom static files (such as style sheets) here,
# relative to this directory. They are copied after the builtin static files,
# so a file named "default.css" will overwrite the builtin "default.css".
html_static_path = ['_static']
htmlhelp_basename = 'Glang Blogs'

formats = ["htmlzip", "pdf", "epub"]

# LaTeX 配置
latex_elements = {
    'papersize': 'a4paper',
    'pointsize': '10pt',
    'preamble': '',
    'figure_align': 'htbp',
}

latex_documents = [
    ('index', 'mkdocs.tex', u'《Golang笔记》',
     u'郑攀', 'manual',),
]

man_pages = [
    ('index', 'pansblog', 'Pan\'s Blog Documentation',
     [u'郑攀'], 1)
]


texinfo_documents = [
    ('index', 'PansBlog', '《Golang博客》',
     u'郑攀', 'PansBlog', '《Golang博客》',
     'Miscellaneous'),
]

on_rtd = os.environ.get('READTHEDOCS', None) == 'True'

if not on_rtd:
    html_theme = 'sphinx_rtd_theme'
    html_theme_path = [sphinx_rtd_theme.get_html_theme_path()]

highlight_langeuage="go,javascript,html"

_exts = "../exts"
sys.path.append(os.path.abspath(_exts))