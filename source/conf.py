import os
import sys
import sphinx_rtd_theme
import platform
# import os
# import sys
# sys.path.insert(0, os.path.abspath('.'))

_exts = "./exts"
sys.path.append(os.path.abspath(_exts))


# -- Project information -----------------------------------------------------

project = 'Rust BLOG'
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

extensions = [
            'sphinx_copybutton',
            'sphinx_markdown_tables',
            'sphinxcontrib.inkscapeconverter',
            'sphinx.ext.autodoc',
            'sphinx.ext.napoleon',
            'sphinx.ext.viewcode',
            'sphinx.ext.autosectionlabel',
            #   'sphinx_simplepdf',
            'chinese_search', 
            'myst_parser', 
            'sphinx.ext.todo',
            ]

# If true, `todo` and `todoList` produce output, else they produce nothing.
todo_include_todos = True


# LaTeX配置
latex_engine = 'xelatex'  # 或者 'pdflatex'，根据你的需求选择

# 根据操作系统选择字体
if platform.system() == 'Windows':
    cjk_font = 'SimSun'
elif platform.system() == 'Darwin':  # macOS
    cjk_font = 'Songti SC'
else:  # Linux
    cjk_font = 'Noto Sans CJK SC'
    
latex_elements = {
    'papersize': 'a4paper',
    'pointsize': '16pt',
    'figure_align': 'htbp',
    'preamble': r'''
    \usepackage{xeCJK}
    \setCJKmainfont{''' + cjk_font + r'''}
    '''
    
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

# Add any paths that contain custom static files (such as style sheets) here,
# relative to this directory. They are copied after the builtin static files,
# so a file named "default.css" will overwrite the builtin "default.css".
html_static_path = ['_static']
html_css_files = [
    'css/custom.css',
]

htmlhelp_basename = 'Glang Blogs'

formats = ["htmlzip", "pdf", "epub"]


latex_documents = [
    ('index', 'mkdocs.tex', u'《Rust笔记》',
     u'郑攀', 'manual',),
]

man_pages = [
    ('index', 'pan blog', 'Pan\'s Blog',
     [u'郑攀'], 1)
]


texinfo_documents = [
    ('index', 'PansBlog', '《Rust博客》',
     u'郑攀', 'PansBlog', '《Rust博客》',
     'Miscellaneous'),
]

on_rtd = os.environ.get('READTHEDOCS', None) == 'True'

if not on_rtd:
    html_theme = 'sphinx_rtd_theme'


highlight_language = "go,javascript,html"



copybutton_prompt_text = "$ "
copybutton_prompt_is_regexp = False
# 给每个 label 加文件路径前缀，避免重复
autosectionlabel_prefix_document = True
autosectionlabel_maxdepth = 1


numfig = True
numfig_secnum_depth = 2

numfig_format = {
    'figure': '图 %s',
    'table': '表 %s',
    'code-block': '代码 %s',
    'section': '节 %s',
}